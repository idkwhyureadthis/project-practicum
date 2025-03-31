-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    phone_number TEXT UNIQUE NOT NULL,
    crypted_password TEXT NOT NULL,
    name TEXT NOT NULL,
    mail TEXT NOT NULL,
    birthday DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;