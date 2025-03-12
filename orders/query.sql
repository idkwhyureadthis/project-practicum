-- name: GetUser :one
SELECT * FROM users
WHERE phone_number = $1;