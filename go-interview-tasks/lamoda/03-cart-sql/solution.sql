-- Решение с собеса Lamoda (исправлено: алиасы, синтаксис SUM, GROUP BY).

-- 1. Кастомеры и элементы корзины
SELECT c.id, c.email, ci.title, ci.amount
FROM customer AS c
LEFT JOIN cart_item AS ci ON c.id = ci.customer_id;

-- 2. Топ-10 покупателей по сумме в корзине
SELECT c.id, c.email, COALESCE(SUM(ci.amount * ci.price), 0) AS total_sum
FROM customer AS c
LEFT JOIN cart_item AS ci ON c.id = ci.customer_id
GROUP BY c.id, c.email
ORDER BY total_sum DESC
LIMIT 10;
