-- As-of date: 2025-06-15

INSERT INTO warehouses (id, name, city) VALUES
    (1, 'Коледино', 'Подольск'),
    (2, 'Казань', 'Казань');

SELECT setval(pg_get_serial_sequence('warehouses', 'id'), 2);

INSERT INTO products (id, sku, name) VALUES
    (1, 'BOOT-42', 'Ботинки зимние'),
    (2, 'COAT-1', 'Пальто'),
    (3, 'HAT-7', 'Шапка');

SELECT setval(pg_get_serial_sequence('products', 'id'), 3);

INSERT INTO warehouse_stock (warehouse_id, product_id, quantity, reserved) VALUES
    (1, 1, 10, 8),   -- BOOT: free=2, demand будет 5+4=9 > 2 → дефицит
    (1, 2, 50, 10),
    (1, 3, 100, 0),
    (2, 1, 5, 5),    -- BOOT: free=0, demand 0 на этом складе
    (2, 2, 20, 0);

INSERT INTO orders (id, warehouse_id, status, created_at) VALUES
    (1001, 1, 'confirmed', '2025-06-14 10:00'),
    (1002, 1, 'pending', '2025-06-14 12:00'),
    (1003, 1, 'shipped', '2025-06-13 09:00'),
    (1004, 1, 'cancelled', '2025-06-12 15:00'),
    (1005, 2, 'confirmed', '2025-06-10 11:00');

SELECT setval(pg_get_serial_sequence('orders', 'id'), 1005);

INSERT INTO order_items (order_id, product_id, qty, price_kopecks) VALUES
    (1001, 1, 5, 499900),   -- BOOT demand
    (1002, 1, 4, 499900),   -- BOOT demand (тот же склад)
    (1001, 2, 1, 1299900),
    (1003, 2, 2, 1299900),
    (1003, 3, 1, 199900),
    (1004, 2, 10, 1299900), -- cancelled — не в спросе
    (1005, 2, 3, 1299900);
