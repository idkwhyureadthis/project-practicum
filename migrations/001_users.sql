-- +goose Up
CREATE TABLE users(
    phone_number TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    about TEXT,
    birthday DATETIME NOT NULL,
    created_at TIMESTAMP
);

-- +goose Down
DROP TABLE users;