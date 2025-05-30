-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    displayed_id INTEGER NOT NULL,
    restaurant_id UUID REFERENCES restaurants(id) NOT NULL,
    total_price FLOAT NOT NULL,
    status TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id)
);


-- +goose Down
DROP TABLE orders;