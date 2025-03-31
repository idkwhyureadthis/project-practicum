-- +goose Up
CREATE TABLE banned_items(
    restaurant_id UUID REFERENCES restaurants(id),
    item_id UUID REFERENCES items(id)
);

-- +goose Down
DROP TABLE banned_items;