CREATE TABLE customer
(
    id      INTEGER PRIMARY KEY,
    email   VARCHAR(100)    NOT NULL,
    country CHAR(2)         NOT NULL
);

CREATE TABLE cart_item
(
    id          INTEGER PRIMARY KEY,
    customer_id INTEGER         NOT NULL,
    title       VARCHAR(20)     NOT NULL,
    amount      INTEGER         NOT NULL,
    price       INTEGER         NOT NULL
);
