-- Magnit-like: loyalty, stores, purchases (PostgreSQL)

CREATE TABLE customers (
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    phone   TEXT NOT NULL UNIQUE
);

CREATE TABLE loyalty_cards (
    id              BIGSERIAL PRIMARY KEY,
    customer_id     BIGINT NOT NULL UNIQUE REFERENCES customers (id),
    card_number     TEXT NOT NULL UNIQUE,
    tier            TEXT NOT NULL CHECK (tier IN ('base', 'silver', 'gold')),
    points_balance  BIGINT NOT NULL DEFAULT 0 CHECK (points_balance >= 0)
);

CREATE TABLE stores (
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    city    TEXT NOT NULL,
    format  TEXT NOT NULL CHECK (format IN ('magnit', 'cosmetic', 'pharmacy'))
);

CREATE TABLE purchases (
    id              BIGSERIAL PRIMARY KEY,
    customer_id     BIGINT NOT NULL REFERENCES customers (id),
    store_id        BIGINT NOT NULL REFERENCES stores (id),
    total_kopecks   BIGINT NOT NULL CHECK (total_kopecks > 0),
    purchased_at    TIMESTAMPTZ NOT NULL
);

CREATE TABLE point_ledger (
    id              BIGSERIAL PRIMARY KEY,
    customer_id     BIGINT NOT NULL REFERENCES customers (id),
    purchase_id     BIGINT REFERENCES purchases (id),
    delta           BIGINT NOT NULL,
    reason          TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_purchases_customer_date ON purchases (customer_id, purchased_at);
CREATE INDEX idx_purchases_store_date ON purchases (store_id, purchased_at);
CREATE INDEX idx_ledger_customer ON point_ledger (customer_id);
