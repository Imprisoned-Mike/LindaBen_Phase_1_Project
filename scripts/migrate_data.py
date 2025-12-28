import ast
import json
import sqlite3
import os

# Configuration
DATA_FILE = '/Users/nativeongfuel/LindaBen_Phase_1_Project/Migration/DATA'
DB_FILE = 'mydatabase.db'

def setup_database(cursor):
    """Creates the necessary tables for the migration."""
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS deliveries (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            contract TEXT,
            school_id INTEGER,
            box_type TEXT,
            notes TEXT,
            scheduled_to DATETIME
        )
    ''')
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            delivery_id INTEGER,
            item TEXT,
            quantity INTEGER,
            vendor_id INTEGER,
            is_internal BOOLEAN,
            status TEXT,
            FOREIGN KEY(delivery_id) REFERENCES deliveries(id)
        )
    ''')

def migrate():
    if not os.path.exists(DATA_FILE):
        print(f"Error: Data file not found at {DATA_FILE}")
        return

    conn = sqlite3.connect(DB_FILE)
    cursor = conn.cursor()
    setup_database(cursor)

    print("Starting migration...")
    count = 0
    
    # Option 1: Handle JSON file if extension matches
    if DATA_FILE.endswith('.json'):
        with open(DATA_FILE, 'r') as f:
            data_list = json.load(f)
            for data in data_list:
                insert_delivery(cursor, data)
                count += 1
        conn.commit()
        conn.close()
        print(f"Migration complete. {count} deliveries imported from JSON.")
        return

    # Option 2: Handle legacy text format
    with open(DATA_FILE, 'r') as f:
        for line in f:
            line = line.strip()
            # Parse only lines starting with create_delivery
            if line.startswith("create_delivery(") and line.endswith(")"):
                try:
                    # Extract the dictionary content: create_delivery({...}) -> {...}
                    content = line[16:-1] 
                    data = ast.literal_eval(content)
                    
                    insert_delivery(cursor, data)
                    count += 1
                except Exception as e:
                    print(f"Failed to parse line: {line[:30]}... Error: {e}")

    conn.commit()
    conn.close()
    print(f"Migration complete. {count} deliveries imported into {DB_FILE}.")

def insert_delivery(cursor, data):
    """Helper function to insert a single delivery and its orders."""
    cursor.execute('''
        INSERT INTO deliveries (contract, school_id, box_type, notes, scheduled_to, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    ''', (
        data.get('contract'),
        data.get('schoolId'),
        data.get('packageType'),
        data.get('notes', ''),
        data.get('scheduledAt')
    ))
    
    delivery_id = cursor.lastrowid

    orders = data.get('orders', [])
    for order in orders:
        # Handle vendor_id -1 (internal) -> NULL to match Go app logic
        vendor_id = order.get('vendorId')
        if vendor_id == -1:
            vendor_id = None

        cursor.execute('''
            INSERT INTO orders (delivery_id, item, quantity, vendor_id, is_internal, status, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        ''', (
            delivery_id,
            order.get('item'),
            order.get('quantity'),
            vendor_id,
            order.get('isInternal'),
            order.get('status')
        ))

if __name__ == "__main__":
    migrate()