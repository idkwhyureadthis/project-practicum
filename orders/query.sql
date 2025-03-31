-- name: LogIn :one
SELECT * FROM users
WHERE phone_number = $1 AND crypted_password = $2;

-- name: CreateUser :one
INSERT INTO users (phone_number, name, crypted_password, mail)
VALUES ($1, $2, $3, $4)
RETURNING *;