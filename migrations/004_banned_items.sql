-- +goose Up
CREATE TABLE banned_items(
    restaurant_id UUID REFERENCES restaurants(id) NOT NULL,
    item_id UUID REFERENCES items(id) NOT NULL
);

-- +goose Down
DROP TABLE banned_items;