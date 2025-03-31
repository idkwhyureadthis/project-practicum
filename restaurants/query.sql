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