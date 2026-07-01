INSERT INTO customers (id, name, phone) VALUES
    (1, 'Светлана', '+79001110001'),
    (2, 'Дмитрий', '+79001110002'),
    (3, 'Елена', '+79001110003'),
    (4, 'Кирилл', '+79001110004'),
    (5, 'Наталья', '+79001110005');

SELECT setval(pg_get_serial_sequence('customers', 'id'), 5);

INSERT INTO loyalty_cards (customer_id, card_number, tier, points_balance) VALUES
    (1, 'MAG-0001', 'gold', 12000),
    (2, 'MAG-0002', 'silver', 8000),
    (3, 'MAG-0003', 'base', 6000),
    (4, 'MAG-0004', 'silver', 1500),
    (5, 'MAG-0005', 'gold', 9000);

INSERT INTO stores (id, name, city, format) VALUES
    (1, 'Магнит №12', 'Краснодар', 'magnit'),
    (2, 'Магнит Косметик', 'Краснодар', 'cosmetic'),
    (3, 'Аптека', 'Сочи', 'pharmacy'),
    (4, 'Магнит №44', 'Ростов-на-Дону', 'magnit');

SELECT setval(pg_get_serial_sequence('stores', 'id'), 4);

INSERT INTO purchases (id, customer_id, store_id, total_kopecks, purchased_at) VALUES
    (101, 1, 1, 450000, '2025-06-10 18:00'),
    (102, 1, 2, 120000, '2025-05-20 12:00'),
    (103, 2, 1, 890000, '2025-06-14 09:00'),
    (104, 2, 4, 310000, '2025-04-01 10:00'),
    (105, 3, 3, 250000, '2025-03-01 11:00'),  -- спящая: last >60d ago
    (106, 4, 1, 99000, '2025-06-01 15:00'),
    (107, 5, 2, 560000, '2025-06-12 20:00'),
    (108, 5, 1, 220000, '2025-06-08 14:00');

SELECT setval(pg_get_serial_sequence('purchases', 'id'), 108);

INSERT INTO point_ledger (customer_id, purchase_id, delta, reason, created_at) VALUES
    (1, 101, 4500, 'accrual_1pct', '2025-06-10'),
    (1, 102, 1200, 'accrual_1pct', '2025-05-20'),
    (2, 103, 8900, 'accrual_1pct', '2025-06-14'),
    (3, 105, 2500, 'accrual_1pct', '2025-03-01');
