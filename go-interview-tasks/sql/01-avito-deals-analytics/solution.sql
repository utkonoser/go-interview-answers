-- Эталонные решения. Для проверки: as-of date = 2025-06-15

-- 1.1 Топ-5 продавцов Москвы по выручке за 30 дней
SELECT
    u.id AS seller_id,
    u.name AS seller_name,
    COUNT(*) AS deals_count,
    SUM(d.price_kopecks) AS revenue_kopecks
FROM deals d
JOIN listings l ON l.id = d.listing_id
JOIN users u ON u.id = l.seller_id
WHERE l.city = 'Москва'
  AND d.completed_at >= DATE '2025-06-15' - INTERVAL '30 days'
  AND d.completed_at < DATE '2025-06-15' + INTERVAL '1 day'
GROUP BY u.id, u.name
ORDER BY revenue_kopecks DESC, seller_id
LIMIT 5;

-- 1.2 Active listings: ≥50 views, zero deals ever
SELECT
    l.id AS listing_id,
    l.title,
    COUNT(v.id) AS views_count,
    (DATE '2025-06-15' - l.created_at::date) AS days_on_site
FROM listings l
JOIN listing_views v ON v.listing_id = l.id
WHERE l.status = 'active'
GROUP BY l.id, l.title, l.created_at
HAVING COUNT(v.id) >= 50
   AND NOT EXISTS (
       SELECT 1 FROM deals d WHERE d.listing_id = l.id
   )
ORDER BY listing_id;

-- 1.3 Воронка по дням (последние 7 дней включая as-of)
WITH days AS (
    SELECT generate_series(
        DATE '2025-06-15' - 6,
        DATE '2025-06-15',
        INTERVAL '1 day'
    )::date AS day
),
views_by_day AS (
    SELECT
        v.viewed_at::date AS day,
        COUNT(DISTINCT v.listing_id) AS viewed_listings
    FROM listing_views v
    WHERE v.viewed_at::date BETWEEN DATE '2025-06-15' - 6 AND DATE '2025-06-15'
    GROUP BY 1
),
deals_by_day AS (
    SELECT
        d.completed_at::date AS day,
        COUNT(*) AS deals_count
    FROM deals d
    WHERE d.completed_at::date BETWEEN DATE '2025-06-15' - 6 AND DATE '2025-06-15'
    GROUP BY 1
)
SELECT
    d.day,
    COALESCE(v.viewed_listings, 0) AS viewed_listings,
    COALESCE(db.deals_count, 0) AS deals_count
FROM days d
LEFT JOIN views_by_day v ON v.day = d.day
LEFT JOIN deals_by_day db ON db.day = d.day
ORDER BY d.day;

-- 1.4 Индекс для ленты (обсуждение, не DML)
-- CREATE INDEX idx_listings_city_status_created
--     ON listings (city, status, created_at DESC)
-- WHERE status = 'active';  -- partial index, если archived/sold много
--
-- Порядок: equality filters (city, status) слева, sort key (created_at DESC) справа.
-- Index-only scan возможен, если все колонки SELECT в индексе или после heap fetch.
