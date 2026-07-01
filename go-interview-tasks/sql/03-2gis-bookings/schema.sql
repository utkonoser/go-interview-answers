-- 2GIS-like: places, rooms, bookings (PostgreSQL)

CREATE TABLE places (
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    city    TEXT NOT NULL
);

CREATE TABLE rooms (
    id          BIGSERIAL PRIMARY KEY,
    place_id    BIGINT NOT NULL REFERENCES places (id),
    name        TEXT NOT NULL,
    capacity    INT NOT NULL CHECK (capacity > 0)
);

CREATE TABLE bookings (
    id          BIGSERIAL PRIMARY KEY,
    room_id     BIGINT NOT NULL REFERENCES rooms (id),
    guest_name  TEXT NOT NULL,
    check_in    DATE NOT NULL,
    check_out   DATE NOT NULL CHECK (check_out > check_in),
    status      TEXT NOT NULL CHECK (status IN ('confirmed', 'cancelled')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_bookings_room_dates ON bookings (room_id, check_in, check_out)
    WHERE status = 'confirmed';

-- Эталон constraint — в design.md (добавляется после миграции данных)
