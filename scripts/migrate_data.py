import ast
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
            contract TEXT,
            school_id INTEGER,
            package_type TEXT,
            notes TEXT,
            scheduled_at DATETIME
        )
    ''')
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
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
    
    with open(DATA_FILE, 'r') as f:
        for line in f:
            line = line.strip()
            # Parse only lines starting with create_delivery
            if line.startswith("create_delivery(") and line.endswith(")"):
                try:
                    # Extract the dictionary content: create_delivery({...}) -> {...}
                    content = line[16:-1] 
                    data = ast.literal_eval(content)

                    # Insert Delivery
                    cursor.execute('''
                        INSERT INTO deliveries (contract, school_id, package_type, notes, scheduled_at)
                        VALUES (?, ?, ?, ?, ?)
                    ''', (
                        data.get('contract'),
                        data.get('schoolId'),
                        data.get('packageType'),
                        data.get('notes', ''),
                        data.get('scheduledAt')
                    ))
                    
                    delivery_id = cursor.lastrowid

                    # Insert Orders
                    orders = data.get('orders', [])
                    for order in orders:
                        cursor.execute('''
                            INSERT INTO orders (delivery_id, item, quantity, vendor_id, is_internal, status)
                            VALUES (?, ?, ?, ?, ?, ?)
                        ''', (
                            delivery_id,
                            order.get('item'),
                            order.get('quantity'),
                            order.get('vendorId'),
                            order.get('isInternal'),
                            order.get('status')
                        ))
                    
                    count += 1
                except Exception as e:
                    print(f"Failed to parse line: {line[:30]}... Error: {e}")

    conn.commit()
    conn.close()
    print(f"Migration complete. {count} deliveries imported into {DB_FILE}.")

if __name__ == "__main__":
    migrate()