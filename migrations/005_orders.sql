-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    displayed_id INTEGER NOT NULL,
    restaurant_id UUID REFERENCES restaurants(id),
    total_price FLOAT NOT NULL,
    status TEXT NOT NULL
);


-- +goose Down
DROP TABLE orders;