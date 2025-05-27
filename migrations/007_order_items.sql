-- +goose Up
CREATE TABLE order_items(
    order_id UUID NOT NULL REFERENCES orders(id),
    item_id UUID NOT NULL REFERENCES items(id)
);


-- +goose Down
DROP TABLE order_items;