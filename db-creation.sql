create table users (
    id integer primary key autoincrement
    , username text not null unique
    , password text not null
    , email text unique
    , isEnabled bool
);

-- Table with constraint check
CREATE TABLE products (
    product_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    price REAL CHECK(price > 0),
    category TEXT DEFAULT 'uncategorized',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table with foreign key
CREATE TABLE orders (
    order_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(product_id)
);

-- Seed data for users
INSERT INTO users (username, password, email, isEnabled) VALUES
('admin', 'admin123', 'admin@example.com', 1),
('alice', 'password456', 'alice@example.com', 1),
('bob', 'securepass789', 'bob@example.com', 1),
('charlie', 'ch@rl!e2020', 'charlie@example.com', 1),
('diana', 'diana_secret', 'diana@example.com', 0);

-- Seed data for products
INSERT INTO products (product_id, name, price, category, created_at) VALUES
(1, 'Laptop Pro', 1299.99, 'electronics', datetime('now')),
(2, 'Wireless Mouse', 49.99, 'electronics', datetime('now')),
(3, 'USB-C Cable', 12.99, 'accessories', datetime('now')),
(4, 'USB Hub', 34.50, 'accessories', datetime('now')),
(5, 'Monitor 27 inch', 349.99, 'electronics', datetime('now')),
(6, 'Keyboard Mechanical', 129.99, 'electronics', datetime('now')),
(7, 'Desk Lamp', 45.00, 'office', datetime('now')),
(8, 'Notebook Set', 14.99, 'office', datetime('now'));

-- Seed data for orders
INSERT INTO orders (user_id, product_id, quantity) VALUES
(1, 1, 1),
(1, 2, 2),
(2, 3, 5),
(2, 4, 1),
(3, 5, 1),
(3, 6, 1),
(4, 7, 3),
(4, 8, 2);
