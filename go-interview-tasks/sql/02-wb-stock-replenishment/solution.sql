-- Эталонные решения. As-of date = 2025-06-15

-- 2.1 Дефицит: demand > free на складе
WITH demand AS (
    SELECT
        o.warehouse_id,
        oi.product_id,
        SUM(oi.qty) AS demand_qty
    FROM orders o
    JOIN order_items oi ON oi.order_id = o.id
    WHERE o.status IN ('pending', 'confirmed')
    GROUP BY o.warehouse_id, oi.product_id
)
SELECT
    p.sku,
    p.name AS product_name,
    w.name AS warehouse_name,
    ws.quantity - ws.reserved AS free_qty,
    d.demand_qty
FROM demand d
JOIN warehouse_stock ws
    ON ws.warehouse_id = d.warehouse_id AND ws.product_id = d.product_id
JOIN products p ON p.id = d.product_id
JOIN warehouses w ON w.id = d.warehouse_id
WHERE d.demand_qty > ws.quantity - ws.reserved
ORDER BY sku, warehouse_name;

-- 2.2 Отчёт по складам за 7 дней
SELECT
    w.name AS warehouse_name,
    COUNT(DISTINCT o.id) AS orders_count,
    SUM(oi.qty * oi.price_kopecks) AS revenue_kopecks
FROM warehouses w
JOIN orders o ON o.warehouse_id = w.id
JOIN order_items oi ON oi.order_id = o.id
WHERE o.status <> 'cancelled'
  AND o.created_at::date BETWEEN DATE '2025-06-15' - 6 AND DATE '2025-06-15'
GROUP BY w.id, w.name
ORDER BY warehouse_name;
