-- +goose Up
CREATE TABLE admins(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT UNIQUE NOT NULL,
    crypted_password TEXT NOT NULL,
    is_superadmin BOOLEAN NOT NULL DEFAULT FALSE,
    restaurant_id UUID REFERENCES restaurants(id),
    crypted_refresh TEXT
);

-- +goose Down
DROP TABLE admins;