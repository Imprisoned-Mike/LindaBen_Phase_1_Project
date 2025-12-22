import httpx
from typing import List, Optional, Literal
from pydantic import BaseModel, Field
import json
import csv
import os
from openai import OpenAI
import difflib
import time
import math
from openai import RateLimitError, APIError 


client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

class Order(BaseModel):
    item: str
    quantity: int = Field(..., description="The quantity of items ordered. ")

    vendorId: int = Field(
        ...,
        description="The unique identifier for the vendor. -1 when vendor is LindaBen.",
    )
    isInternal: bool = Field(..., description="True when vendor is LindaBen.")

    # Conditional fields based on TS comments
    completedAt: Optional[str] = Field(
        None, description="ISO timestamp when the order was completed."
    )
    status: Literal["pending", "confirmed", "completed", "cancelled"] = Field(
        ..., description="The current status of the order."
    )


class Delivery(BaseModel):
    orders: List[Order]
    contract: Literal["hold", "completed"]
    schoolId: int

    packageType: Literal["standard", "premium", "standard-holiday", "premium-holiday"]
    notes: Optional[str] = Field(
        None, description="Leave empty unless there are special instructions."
    )

    scheduledAt: Optional[str] = Field(
        None,
        description="ISO timestamp for when the delivery is scheduled.",
    )


SYSTEM = """
You are a helpful agent that converts a human written table into fixed data models in JSON format. A single row of the table looks like this:
School Name, Scheduled date, Address, Quantity, Scheduled Time, Package Type, Contract Status, none, [...Orders columns], none, none, Notes
For an order column, the content refers to the vendor name, if empty, then there is no order with such item. 
The order items are (sequentially): "Sm Produce" "Lg Produce" "Rx Box" "Chicken" "Turkey" "Hallal Meat" "Eggs" "Fresh Milk" "Bread"
Each row is converted into one Delivery object, multiple Order objects, schools and vendors already exist.
You will process the conversion in a multi-turn way, each time calling a tool I provided with you, and process its response. 

Assume year is 2026 for all dates.

Call finish_row (or emit a clear “row done” signal). Only call abort on unrecoverable errors. Only call finalize batch after all rows in the last batch are completed.

NOTE: LindaBen is the internal vendor. No need to search for it. When searching for schools, search with minimum keywords, "Hollywood" instead of "Hollywood ES (70, 8am)"
"""

TOOLS = [
    {
        "type": "function",
        "function": {
            "name": "search_schools",
            "description": "Search for schools by name.",
            "parameters": {
                "type": "object",
                "properties": {
                    "name": {
                        "type": "string",
                        "description": "The name to search",
                    },
                },
                "required": ["name"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "search_vendors",
            "description": "Search for vendors by name.",
            "parameters": {
                "type": "object",
                "properties": {
                    "name": {
                        "type": "string",
                        "description": "The name to search",
                    },
                },
                "required": ["name"],
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "create_delivery",
            "description": "Create a delivery with orders.",
            "parameters": Delivery.model_json_schema(),
        },
    },
    {
        "type": "function",
        "function": {
            "name": "abort",
            "description": "Abort the conversion process with a reason.",
            "parameters": {
                "type": "object",
                "properties": {
                    "reason": {
                        "type": "string",
                        "description": "The reason for aborting the process.",
                    },
                    "is_error": {
                        "type": "boolean",
                    },
                },
                "required": ["reason", "is_error"],
            },
        },
    },{
        "type": "function",
        "function": {
            "name": "finalize_batch",
            "description": "Signal that all rows in the final batch were processed; persist any mappings.",
            "parameters": {"type": "object", "properties": {}},
        },
    }
]

MODEL = "gpt-4.1"

DATA = """
School Name, ScheduledAt, Address, Quantity, ScheduledTime, PackageType, ContractStatus, Empty, Sm Produce, Lg Produce, Rx Box, Chicken, Turkey, Hallal Meat, Eggs, Fresh Milk, Bread, Empty, Empty, Notes
"""

# parse and convert to friendly format
data_lines = DATA.strip().split("\n")
data_reader = csv.reader(data_lines)
data_rows = list(data_reader)
header, all_rows = data_rows[0], data_rows[1:]
# print(json.dumps(data_rows, indent=2))
# exit(0)

messages = [
    {
        "role": "system",
        "content": SYSTEM,
    },
    {
        "role": "user",
        "content": json.dumps(data_rows),
    },
]

SCHOOLS_FILE = "schools.json"
VENDORS_FILE = "vendors.json"

# Load existing mappings if they exist
if os.path.exists(SCHOOLS_FILE):
    with open(SCHOOLS_FILE, "r") as f:
        schools = json.load(f)
else:
    schools = {}

if os.path.exists(VENDORS_FILE):
    with open(VENDORS_FILE, "r") as f:
        vendors = json.load(f)
else:
    vendors = {}

def handle_tool_call(tool_call):
    global schools, vendors
    function_name = tool_call.function.name
    args = json.loads(tool_call.function.arguments)
    
    if function_name == "search_schools":
        name_query = args["name"]

        for name, sid in schools.items():
            if name == name_query:
                return json.dumps({"id": int(sid)})
        matches = difflib.get_close_matches(name_query, schools.keys(), n=1, cutoff=0.5)
        if matches:
            best = matches[0]
            return json.dumps({"id": int(schools[best])})

        school_id = 1 + max([int(v) for v in schools.values()] + [0])
        schools[name_query] = school_id
        return json.dumps({"id": school_id})
       # for name, sid in schools.items():
       #     if name == args["name"]:
       #         return json.dumps({"id": int(sid)})
       # school_id = 1 + max([int(k) for k in schools.keys()] + [0])
       # schools[str(school_id)] = args["name"]
       # return json.dumps({"id": school_id})
    
    elif function_name == "search_vendors":
        name_query = args["name"]

        for name, vid in vendors.items():
            if name == name_query:
                return json.dumps({"id": int(vid)})
        matches = difflib.get_close_matches(name_query, vendors.keys(), n=1, cutoff=0.5)
        if matches:
            best = matches[0]
            return json.dumps({"id": int(vendors[best])})
        vendor_id = 101 + max([int(v) for v in vendors.values()] + [0])
        vendors[name_query] = vendor_id
        return json.dumps({"id": vendor_id})
        #for name, vid in vendors.items():
        #    if name == args["name"]:
        #        return json.dumps({"id": int(vid)})
        #vendor_id = 101 + max([int(k) for k in vendors.keys()] + [0])
        #vendors[str(vendor_id)] = args["name"]
        #return json.dumps({"id": vendor_id})
    elif function_name == "create_delivery":
        print("Delivery created:", args)
        return json.dumps({"delivery_created": True})
    
    elif function_name == "abort":
        print("Aborting:", args["reason"])
        exit(0)
    elif function_name == "finalize_batch":
        with open(SCHOOLS_FILE, "w") as f:
            json.dump(schools, f, indent=2)
        with open(VENDORS_FILE, "w") as f:
            json.dump(vendors, f, indent=2)
        return json.dumps({"batch_finalized": True})

def finalize_batch():
    with open(SCHOOLS_FILE, "w") as f:
        json.dump(schools, f, indent=2)
    with open(VENDORS_FILE, "w") as f:
        json.dump(vendors, f, indent=2)
    return json.dumps({"batch_finalized": True})

# ----- Helpers: chunking and retry -----
def chunked(seq, size):
    for i in range(0, len(seq), size):
        yield seq[i:i + size]

MAX_RETRIES = 6
BACKOFF = 1.6
INITIAL_DELAY = 1.0
PER_CALL_PAUSE = 0.25  # throttle between tool-turns
PER_BATCH_PAUSE = 10  # throttle between batches

def chat_with_retry(**kwargs):
    delay = INITIAL_DELAY
    for attempt in range(MAX_RETRIES):
        try:
            return client.chat.completions.create(**kwargs)
        except RateLimitError:
            time.sleep(delay)
            delay *= BACKOFF
        except APIError as e:
            status = getattr(e, "status", None)
            if status in (429, 500, 502, 503):
                time.sleep(delay)
                delay *= BACKOFF
            else:
                raise
    raise RuntimeError("Exceeded retries on chat.completions.create")

# ----- Batched processing -----
BATCH_SIZE = 5
max_iterations = 500  # per-batch tool-calling loop
i=0

total_batches = math.ceil(len(all_rows) / BATCH_SIZE)

for batch_idx, batch_rows in enumerate(chunked(all_rows, BATCH_SIZE), start=1):
    is_last_batch = batch_idx == total_batches
    # Only send header + this batch so the prompt stays small.
    batch_table = [header] + batch_rows
    system_content = SYSTEM + (
        "\nThis is the last batch. After finishing all rows, call finalize_batch."
        if is_last_batch else
        "\nDo not call finalize_batch in this batch; only call finish_row after each row."
    )

    messages = [
        {"role": "system", "content": system_content},
        {"role": "user", "content": json.dumps(batch_table)},
    ]

    response = chat_with_retry(
        model=MODEL,
        messages=messages,
        tools=TOOLS,
        tool_choice="auto",
    )

    while getattr(response.choices[0].message, "tool_calls", None) and i < max_iterations:
       i += 1
       messages.append(response.choices[0].message)

       print(f"Batch {batch_idx}/{total_batches}, Iteration {i}: Model called {len(response.choices[0].message.tool_calls)} tool(s)")

       for tool_call in response.choices[0].message.tool_calls:
          function_name = tool_call.function.name
          function_args = json.loads(tool_call.function.arguments or "{}")

          if function_name == "abort":
                print("Aborting conversion:")
                print(f"  Reason: {function_args.get('reason')}")
                raise SystemExit(0)

          print(f"  → {function_name}({function_args})")
          function_response = handle_tool_call(tool_call)
          print(f"    ← {function_response}")

          messages.append({
             "role": "tool",
             "tool_call_id": tool_call.id,
             "name": function_name,
             "content": function_response,
            })

        # Throttle between model turns inside a batch
       time.sleep(PER_CALL_PAUSE)

       response = chat_with_retry(
          model=MODEL,
          messages=messages,
          tools=TOOLS,
          tool_choice="auto",
        )
    # Throttle between batches
    time.sleep(PER_BATCH_PAUSE)
print("All batches processed.")