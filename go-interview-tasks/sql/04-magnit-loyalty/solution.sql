-- Эталонные решения. As-of date = 2025-06-15

-- 4.1 Топ-5 по тратам за 90 дней
SELECT
    c.id AS customer_id,
    c.name AS customer_name,
    SUM(p.total_kopecks) AS spend_kopecks,
    lc.points_balance,
    lc.tier
FROM customers c
JOIN purchases p ON p.customer_id = c.id
JOIN loyalty_cards lc ON lc.customer_id = c.id
WHERE p.purchased_at >= DATE '2025-06-15' - INTERVAL '90 days'
  AND p.purchased_at < DATE '2025-06-15' + INTERVAL '1 day'
GROUP BY c.id, c.name, lc.points_balance, lc.tier
ORDER BY spend_kopecks DESC, customer_id
LIMIT 5;

-- 4.2 Спящие клиенты с баллами
WITH last_purchase AS (
    SELECT customer_id, MAX(purchased_at)::date AS last_purchase_date
    FROM purchases
    GROUP BY customer_id
)
SELECT
    c.id AS customer_id,
    c.name AS customer_name,
    lc.points_balance,
    lp.last_purchase_date,
    (DATE '2025-06-15' - lp.last_purchase_date) AS days_since_purchase
FROM customers c
JOIN loyalty_cards lc ON lc.customer_id = c.id
JOIN last_purchase lp ON lp.customer_id = c.id
WHERE lc.points_balance >= 5000
  AND lp.last_purchase_date < DATE '2025-06-15' - 60
ORDER BY lc.points_balance DESC, customer_id;

-- 4.3 Средний чек по формату магазина за 30 дней
SELECT
    s.format,
    COUNT(*) AS receipts_count,
    ROUND(AVG(p.total_kopecks))::bigint AS avg_check_kopecks
FROM purchases p
JOIN stores s ON s.id = p.store_id
WHERE p.purchased_at::date BETWEEN DATE '2025-06-15' - 29 AND DATE '2025-06-15'
GROUP BY s.format
ORDER BY format;
