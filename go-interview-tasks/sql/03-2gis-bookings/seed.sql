INSERT INTO places (id, name, city) VALUES
    (1, 'Морской бриз', 'Сочи'),
    (2, 'Горный воздух', 'Сочи'),
    (3, 'Невский', 'Санкт-Петербург');

SELECT setval(pg_get_serial_sequence('places', 'id'), 3);

INSERT INTO rooms (id, place_id, name, capacity) VALUES
    (101, 1, 'Стандарт', 2),
    (102, 1, 'Люкс', 4),
    (201, 2, 'Стандарт', 2),
    (301, 3, 'Эконом', 2);

SELECT setval(pg_get_serial_sequence('rooms', 'id'), 301);

INSERT INTO bookings (id, room_id, guest_name, check_in, check_out, status) VALUES
    (1, 101, 'Алиса', '2025-06-15', '2025-06-18', 'confirmed'),  -- 3 ночи в окне occupancy
    (2, 101, 'Борис', '2025-06-18', '2025-06-20', 'confirmed'),  -- стык без overlap
    (3, 102, 'Вика', '2025-06-16', '2025-06-19', 'confirmed'),
    (4, 201, 'Глеб', '2025-06-15', '2025-06-17', 'confirmed'),
    (5, 101, 'Дина', '2025-06-17', '2025-06-19', 'cancelled'),   -- не в отчётах
    (6, 101, 'Егор', '2025-06-17', '2025-06-19', 'confirmed'),   -- overlap с id=1 для 3.2
    (7, 102, 'Жанна', '2025-07-01', '2025-07-04', 'confirmed'), -- Сочи, занят на 3.1–5.07
    (8, 201, 'Игорь', '2025-07-02', '2025-07-05', 'confirmed');

SELECT setval(pg_get_serial_sequence('bookings', 'id'), 8);
