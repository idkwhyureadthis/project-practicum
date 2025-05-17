-- +goose Up

-- +goose StatementBegin
CREATE TABLE admins(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT UNIQUE NOT NULL,
    crypted_password TEXT NOT NULL,
    is_superadmin BOOLEAN NOT NULL DEFAULT FALSE,
    restaurant_id UUID REFERENCES restaurants(id),
    crypted_refresh TEXT
);

INSERT INTO admins(login, crypted_password, is_superadmin)
VALUES ("admin", encode(sha256("12345"::bytea), 'hex'), true)

-- +goose StatementEnd

-- +goose Down
DROP TABLE admins;