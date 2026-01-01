# The non-invasive data migration script

import time
import asyncio
from typing import Dict, List
import httpx
import re
from datetime import datetime
from zoneinfo import ZoneInfo

# Configuration
DATA_FILE = "./DATA"
API_BASE = "https://rx-api.harvey-l.com/api"

transport = httpx.HTTPTransport(retries=3)
async_transport = httpx.AsyncHTTPTransport(retries=3)


def login(email, password):
    response = httpx.post(
        f"{API_BASE}/auth/login", json={"email": email, "password": password}
    )
    response.raise_for_status()
    return response.json()["token"]


token = login("admin@logistics.com", "password")
print("Logged in successfully")


def create_school(school):
    response = httpx.post(
        f"{API_BASE}/schools",
        json=school,
    )
    response.raise_for_status()
    print(f"Created school with ID: {response.json()['id']}")


def create_vendor(vendor):
    response = httpx.post(
        f"{API_BASE}/vendors",
        json=vendor,
    )
    response.raise_for_status()
    print(f"Created vendor with ID: {response.json()['id']}")


async def delete_delivery_async(
    client: httpx.AsyncClient, delivery_id: int, sem: asyncio.Semaphore
):
    async with sem:
        response = await client.delete(f"{API_BASE}/deliveries/{delivery_id}")
        response.raise_for_status()
        print(f"Deleted delivery with ID: {delivery_id}")


async def delete_all_deliveries_parallel(concurrency: int = 20):
    async with httpx.AsyncClient(
        headers={"Authorization": f"Bearer {token}"},
        transport=async_transport,
    ) as client:
        response = await client.get(f"{API_BASE}/deliveries?pageSize=9999")
        response.raise_for_status()
        deliveries = response.json()["data"]

        sem = asyncio.Semaphore(concurrency)
        await asyncio.gather(
            *(delete_delivery_async(client, d["id"], sem) for d in deliveries)
        )


async def create_delivery_async(
    client: httpx.AsyncClient, delivery: Dict, sem: asyncio.Semaphore
):
    async with sem:
        response = await client.post(f"{API_BASE}/deliveries", json=delivery)
        response.raise_for_status()
        print(f"Created delivery with ID: {response.json()['id']}")


async def create_deliveries_parallel(deliveries: List[Dict], concurrency: int = 20):
    async with httpx.AsyncClient(
        headers={"Authorization": f"Bearer {token}"},
        transport=async_transport,
    ) as client:
        sem = asyncio.Semaphore(concurrency)
        await asyncio.gather(
            *(create_delivery_async(client, d, sem) for d in deliveries)
        )


def parse_mappings(line: str) -> Dict[str, int]:
    # line is
    # Vendors = {Miller: 101, Rahll: 202, Albright: 303}
    line = line.split("{")[1].split("}")[0]
    mappings = {}
    for item in line.split(","):
        name, id_ = item.split(":")
        mappings[name.strip()] = int(id_.strip())
    return mappings


def search_osm(name: str) -> Dict:
    base = "https://nominatim.openstreetmap.org/search"
    response = httpx.get(
        base,
        params={"q": name, "format": "json", "limit": 1},
        headers={"User-Agent": "LindaBen Foundation. Reports to hli977@gatech.edu"},
    )
    response.raise_for_status()
    results = response.json()
    if len(results) == 0:
        raise ValueError(f"No OSM results for {name}")

    print(f"OSM search for {name}: {results[0]}")

    return {
        "coordinate": {
            "lng": float(results[0]["lon"]),
            "lat": float(results[0]["lat"]),
        },
        "address": results[0]["display_name"],
    }


# Parse the schools and vendors mapping
vendor_mappings = {}
actual_vendor_mappings = {}
school_mappings = {}
actual_school_mappings = {}
for line in open(DATA_FILE):
    if line.startswith("Vendors"):
        vendor_mappings = parse_mappings(line)
    elif line.startswith("Schools"):
        school_mappings = parse_mappings(line)

print("Vendor mappings:", vendor_mappings)
print("School mappings:", school_mappings)

# Ensure schools and vendors exist
existing_schools = httpx.get(f"{API_BASE}/schools?pageSize=9999").json()["data"]
existing_vendors = httpx.get(f"{API_BASE}/vendors?pageSize=9999").json()["data"]

for school in existing_schools:
    actual_school_mappings[school["name"]] = school["id"]
for vendor in existing_vendors:
    actual_vendor_mappings[vendor["name"]] = vendor["id"]

for school_name, school_id in school_mappings.items():
    if school_name not in actual_school_mappings:
        print(f"Creating school: {school_name}")

        try:
            osm = search_osm(school_name)
            osm["name"] = school_name
            create_school(osm)
        except Exception as e:
            print("No OSM data found, creating with default values.")
            create_school(
                {
                    "name": school_name,
                    "address": "Unknown",
                    "coordinate": {"lat": 0, "lon": 0},
                }
            )

        time.sleep(1)  # Be nice to OSM servers

for vendor_name, vendor_id in vendor_mappings.items():
    if vendor_name not in actual_vendor_mappings:
        print(f"Creating vendor: {vendor_name}")
        create_vendor({"name": vendor_name})

print("All schools and vendors are ensured to exist.")

# create id-to-id map
school_id_map = {
    old_id: actual_school_mappings[school_name]
    for school_name, old_id in school_mappings.items()
}
vendor_id_map = {
    old_id: actual_vendor_mappings[vendor_name]
    for vendor_name, old_id in vendor_mappings.items()
}

deliveries_to_create = []
for line in open(DATA_FILE).readlines():
    if line.startswith("create_delivery"):
        delivery_data = eval(re.search(r"create_delivery\((.*)\)", line).group(1))  # pyright:ignore

        # Re-format time. TZ = EST
        dt_naive = datetime.fromisoformat(delivery_data["scheduledAt"])
        dt_eastern = dt_naive.replace(tzinfo=ZoneInfo("America/New_York"))
        delivery_data["scheduledAt"] = dt_eastern.isoformat()

        is_past = dt_eastern < datetime.now(tz=ZoneInfo("America/New_York"))

        # Re-map
        delivery_data["schoolId"] = school_id_map[delivery_data["schoolId"]]
        for o in delivery_data["orders"]:
            if o["vendorId"] != -1:
                o["vendorId"] = vendor_id_map[o["vendorId"]]
            if is_past:
                # Assume all past deliveries were delivered on time
                o["status"] = "completed"
                o["completedAt"] = dt_eastern.isoformat()

        deliveries_to_create.append(delivery_data)


async def migrate_deliveries():
    await delete_all_deliveries_parallel()
    await create_deliveries_parallel(deliveries_to_create)


asyncio.run(migrate_deliveries())
