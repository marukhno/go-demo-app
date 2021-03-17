CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    ticker varchar(100) NOT NULL,
    price NUMERIC(12, 2),
    created timestamp
);
