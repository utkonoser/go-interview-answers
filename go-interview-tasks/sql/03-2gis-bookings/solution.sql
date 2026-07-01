-- Эталонные решения

-- 3.1 Загрузка отелей: 2025-06-15 .. 2025-06-21 (7 ночей)
WITH period AS (
    SELECT
        DATE '2025-06-15' AS start_day,
        DATE '2025-06-22' AS end_day,  -- exclusive
        7 AS nights_in_period
),
room_nights AS (
    SELECT
        r.place_id,
        b.room_id,
        gs.night::date AS night
    FROM bookings b
    JOIN rooms r ON r.id = b.room_id
    CROSS JOIN period p
    CROSS JOIN LATERAL generate_series(
        GREATEST(b.check_in, p.start_day),
        LEAST(b.check_out, p.end_day) - 1,
        INTERVAL '1 day'
    ) AS gs(night)
    WHERE b.status = 'confirmed'
      AND b.check_in < p.end_day
      AND b.check_out > p.start_day
),
agg AS (
    SELECT place_id, COUNT(*) AS booked_nights
    FROM room_nights
    GROUP BY place_id
)
SELECT
    pl.name AS place_name,
    pl.city,
    COUNT(r.id) AS room_count,
    COALESCE(a.booked_nights, 0) AS booked_nights,
    ROUND(
        100.0 * COALESCE(a.booked_nights, 0)
        / (COUNT(r.id) * (SELECT nights_in_period FROM period)),
        1
    ) AS occupancy_pct
FROM places pl
JOIN rooms r ON r.place_id = pl.id
LEFT JOIN agg a ON a.place_id = pl.id
GROUP BY pl.id, pl.name, pl.city, a.booked_nights
ORDER BY occupancy_pct DESC, place_name;

-- 3.2 Пересекающиеся confirmed-брони одного номера
SELECT
    LEAST(a.id, b.id) AS booking_id_a,
    GREATEST(a.id, b.id) AS booking_id_b,
    a.room_id
FROM bookings a
JOIN bookings b
    ON a.room_id = b.room_id
   AND a.id < b.id
WHERE a.status = 'confirmed'
  AND b.status = 'confirmed'
  AND daterange(a.check_in, a.check_out, '[)') && daterange(b.check_in, b.check_out, '[)')
ORDER BY a.room_id, booking_id_a;

-- 3.3 Свободные номера в Сочи на 2025-07-01 .. 2025-07-05 (4 ночи)
SELECT
    r.id AS room_id,
    r.name AS room_name,
    p.name AS place_name
FROM rooms r
JOIN places p ON p.id = r.place_id
WHERE p.city = 'Сочи'
  AND NOT EXISTS (
      SELECT 1
      FROM bookings b
      WHERE b.room_id = r.id
        AND b.status = 'confirmed'
        AND daterange(b.check_in, b.check_out, '[)') &&
            daterange(DATE '2025-07-01', DATE '2025-07-05', '[)')
  )
ORDER BY place_name, room_id;
