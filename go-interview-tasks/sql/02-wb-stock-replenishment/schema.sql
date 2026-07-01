-- WB / Lamoda-like: warehouses, stock, orders (PostgreSQL)

CREATE TABLE warehouses (
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    city    TEXT NOT NULL
);

CREATE TABLE products (
    id      BIGSERIAL PRIMARY KEY,
    sku     TEXT NOT NULL UNIQUE,
    name    TEXT NOT NULL
);

-- Остатки и резерв на складе (см. design.md)
CREATE TABLE warehouse_stock (
    warehouse_id    BIGINT NOT NULL REFERENCES warehouses (id),
    product_id      BIGINT NOT NULL REFERENCES products (id),
    quantity        INT NOT NULL CHECK (quantity >= 0),
    reserved        INT NOT NULL DEFAULT 0 CHECK (reserved >= 0),
    CHECK (reserved <= quantity),
    PRIMARY KEY (warehouse_id, product_id)
);

CREATE TABLE orders (
    id              BIGSERIAL PRIMARY KEY,
    warehouse_id    BIGINT NOT NULL REFERENCES warehouses (id),
    status          TEXT NOT NULL CHECK (status IN ('pending', 'confirmed', 'shipped', 'cancelled')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT NOT NULL REFERENCES orders (id),
    product_id      BIGINT NOT NULL REFERENCES products (id),
    qty             INT NOT NULL CHECK (qty > 0),
    price_kopecks   BIGINT NOT NULL CHECK (price_kopecks > 0)
);

CREATE INDEX idx_order_items_order ON order_items (order_id);
CREATE INDEX idx_order_items_product ON order_items (product_id);
CREATE INDEX idx_orders_status_created ON orders (status, created_at);
