-- name: GetAdmin :one
SELECT * FROM admins
WHERE login = $1 AND crypted_password = $2;


-- name: UpdateRefresh :exec 
UPDATE admins
SET crypted_refresh = $1 where id = $2;


-- name: GetRefresh :one
SELECT crypted_refresh from admins
WHERE id = $1;


-- name: AddRestaurant :one
INSERT INTO restaurants
(coordinates, name, open_time, close_time) VALUES
($1, $2, $3, $4)
RETURNING id;


-- name: GetRestaurants :many
SELECT * FROM restaurants;


-- name: CreateAdmin :one
INSERT INTO admins (login, crypted_password, restaurant_id)
VALUES ($1, $2, $3)
RETURNING id;

-- name: CreateItem :one
INSERT INTO items (name, description, sizes, prices, photos)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: BanItem :exec
INSERT INTO banned_items
(item_id, restaurant_id)
VALUES ($1, $2);

-- name: UnbanItem :exec
DELETE FROM banned_items
WHERE item_id = $1 AND restaurant_id = $2;

-- name: GetBannedItems :many
SELECT item_id FROM banned_items
WHERE restaurant_id = $1;


-- name: GetAdminRestaurant :one
SELECT restaurant_id from admins
WHERE id = $1;

-- name: GetItems :many
SELECT 
    i.*,
    CASE 
        WHEN bi.item_id IS NULL THEN true 
        ELSE false 
    END AS is_available
FROM items i
LEFT JOIN banned_items bi ON i.id = bi.item_id AND bi.restaurant_id = $1;


-- name: GetRestaurantOrders :many
SELECT * FROM orders
WHERE restaurant_id = $1;

