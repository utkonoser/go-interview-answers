-- Avito-like: listings, views, deals (PostgreSQL)

CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    city       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE listings (
    id              BIGSERIAL PRIMARY KEY,
    seller_id       BIGINT NOT NULL REFERENCES users (id),
    title           TEXT NOT NULL,
    price_kopecks   BIGINT NOT NULL CHECK (price_kopecks > 0),
    city            TEXT NOT NULL,
    status          TEXT NOT NULL CHECK (status IN ('active', 'sold', 'archived')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE listing_views (
    id          BIGSERIAL PRIMARY KEY,
    listing_id  BIGINT NOT NULL REFERENCES listings (id),
    viewer_id   BIGINT REFERENCES users (id),
    viewed_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE deals (
    id              BIGSERIAL PRIMARY KEY,
    listing_id      BIGINT NOT NULL REFERENCES listings (id),
    buyer_id        BIGINT NOT NULL REFERENCES users (id),
    price_kopecks   BIGINT NOT NULL CHECK (price_kopecks > 0),
    completed_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
