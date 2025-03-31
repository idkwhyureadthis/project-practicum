-- name: LogIn :one
SELECT * FROM users
WHERE phone_number = $1 AND crypted_password = $2;

-- name: CreateUser :one
INSERT INTO users (id, phone_number, name, crypted_password, mail)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetRefresh :one
SELECT crypted_refresh FROM users WHERE id = $1;

-- name: UpdateRefresh :exec
UPDATE users SET crypted_refresh = $2 WHERE id = $1;