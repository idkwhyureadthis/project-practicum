-- name: CheckAdmin :one
SELECT COUNT(*) FROM admins WHERE login = 'admin';

-- name: SetupAdmin :exec
INSERT INTO admins
(login, crypted_password, is_superadmin)
VALUES ('admin', $1, TRUE);


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