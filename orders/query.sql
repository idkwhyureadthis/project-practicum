-- name: LogIn :one
SELECT * FROM users
WHERE phone_number = $1 AND crypted_password = $2;

-- name: CreateUser :one
INSERT INTO users (phone_number, name, crypted_password, mail)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT name, id, phone_number, mail, birthday, created_at FROM users WHERE id = $1;

-- name: GetRefresh :one
SELECT crypted_refresh FROM users WHERE id = $1;

-- name: UpdateRefresh :exec
UPDATE users SET crypted_refresh = $2 WHERE id = $1;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 AND user_id = $2;

-- name: GetUserOrders :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY displayed_id;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1 AND user_id = $2;

-- name: GetOrderByDisplayedId :one
SELECT * FROM orders
WHERE displayed_id = $1 AND restaurant_id = $2;

-- name: GetItemById :one
SELECT * FROM items
WHERE id = $1;


-- name: CreateOrder :one
INSERT INTO orders (id, displayed_id, restaurant_id, total_price, status, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, item_id)
VALUES ($1, $2);