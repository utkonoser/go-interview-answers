-- Фиксированная «сегодняшняя» дата для воспроизводимости: 2025-06-15
-- В solution.sql используем DATE '2025-06-15' вместо CURRENT_DATE при проверке.

INSERT INTO users (id, name, city, created_at) VALUES
    (1, 'Иван', 'Москва', '2024-01-01'),
    (2, 'Мария', 'Москва', '2024-02-01'),
    (3, 'Пётр', 'Санкт-Петербург', '2024-03-01'),
    (4, 'Анна', 'Москва', '2024-04-01'),
    (5, 'Олег', 'Казань', '2024-05-01');

SELECT setval(pg_get_serial_sequence('users', 'id'), 5);

INSERT INTO listings (id, seller_id, title, price_kopecks, city, status, created_at) VALUES
    (101, 1, 'iPhone 13', 4500000, 'Москва', 'active', '2025-05-01'),
    (102, 1, 'MacBook Air', 8900000, 'Москва', 'sold', '2025-04-01'),
    (103, 2, 'Диван', 2500000, 'Москва', 'active', '2025-03-15'),
    (104, 2, 'Велосипед', 1200000, 'Москва', 'active', '2025-06-01'),
    (105, 3, 'Книги', 50000, 'Санкт-Петербург', 'active', '2025-05-20'),
    (106, 4, 'Пылесос', 800000, 'Москва', 'archived', '2025-01-10');

SELECT setval(pg_get_serial_sequence('listings', 'id'), 106);

-- 101 и 104 — ≥50 просмотров, без сделок; 103 — ≥50 просмотров, сделка есть
INSERT INTO listing_views (listing_id, viewer_id, viewed_at)
SELECT 101, 4, ts FROM generate_series(
    '2025-06-01'::timestamptz, '2025-06-14'::timestamptz, '6 hours'
) AS ts;

INSERT INTO listing_views (listing_id, viewer_id, viewed_at)
SELECT 103, 4, ts FROM generate_series(
    '2025-05-20'::timestamptz, '2025-06-05'::timestamptz, '8 hours'
) AS ts;

INSERT INTO listing_views (listing_id, viewer_id, viewed_at)
SELECT 104, 1, ts FROM generate_series(
    '2025-06-01'::timestamptz, '2025-06-14'::timestamptz, '6 hours'
) AS ts;

INSERT INTO listing_views (listing_id, viewer_id, viewed_at) VALUES
    (105, 1, '2025-06-13');

INSERT INTO deals (listing_id, buyer_id, price_kopecks, completed_at) VALUES
    (102, 4, 8700000, '2025-06-10'),
    (103, 4, 2400000, '2025-06-01'),
    (105, 2, 45000, '2025-06-08');
