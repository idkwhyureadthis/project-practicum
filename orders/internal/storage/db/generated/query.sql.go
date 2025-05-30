// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package generated

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (id, displayed_id, restaurant_id, total_price, status, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, displayed_id, restaurant_id, total_price, status, user_id
`

type CreateOrderParams struct {
	ID           uuid.UUID `json:"id"`
	DisplayedID  int32     `json:"displayed_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	TotalPrice   float64   `json:"total_price"`
	Status       string    `json:"status"`
	UserID       uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, createOrder,
		arg.ID,
		arg.DisplayedID,
		arg.RestaurantID,
		arg.TotalPrice,
		arg.Status,
		arg.UserID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.DisplayedID,
		&i.RestaurantID,
		&i.TotalPrice,
		&i.Status,
		&i.UserID,
	)
	return i, err
}

const createOrderItem = `-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, item_id)
VALUES ($1, $2)
`

type CreateOrderItemParams struct {
	OrderID uuid.UUID `json:"order_id"`
	ItemID  uuid.UUID `json:"item_id"`
}

func (q *Queries) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) error {
	_, err := q.db.Exec(ctx, createOrderItem, arg.OrderID, arg.ItemID)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (phone_number, name, crypted_password, mail)
VALUES ($1, $2, $3, $4)
RETURNING id, phone_number, crypted_password, name, mail, birthday, crypted_refresh, created_at
`

type CreateUserParams struct {
	PhoneNumber     string `json:"phone_number"`
	Name            string `json:"name"`
	CryptedPassword string `json:"crypted_password"`
	Mail            string `json:"mail"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.PhoneNumber,
		arg.Name,
		arg.CryptedPassword,
		arg.Mail,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.PhoneNumber,
		&i.CryptedPassword,
		&i.Name,
		&i.Mail,
		&i.Birthday,
		&i.CryptedRefresh,
		&i.CreatedAt,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1 AND user_id = $2
`

type DeleteOrderParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteOrder(ctx context.Context, arg DeleteOrderParams) error {
	_, err := q.db.Exec(ctx, deleteOrder, arg.ID, arg.UserID)
	return err
}

const getItemById = `-- name: GetItemById :one
SELECT id, name, description, sizes, prices, photos FROM items
WHERE id = $1
`

func (q *Queries) GetItemById(ctx context.Context, id uuid.UUID) (Item, error) {
	row := q.db.QueryRow(ctx, getItemById, id)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Sizes,
		&i.Prices,
		&i.Photos,
	)
	return i, err
}

const getOrderByDisplayedId = `-- name: GetOrderByDisplayedId :one
SELECT id, displayed_id, restaurant_id, total_price, status, user_id FROM orders
WHERE displayed_id = $1 AND restaurant_id = $2
`

type GetOrderByDisplayedIdParams struct {
	DisplayedID  int32     `json:"displayed_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
}

func (q *Queries) GetOrderByDisplayedId(ctx context.Context, arg GetOrderByDisplayedIdParams) (Order, error) {
	row := q.db.QueryRow(ctx, getOrderByDisplayedId, arg.DisplayedID, arg.RestaurantID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.DisplayedID,
		&i.RestaurantID,
		&i.TotalPrice,
		&i.Status,
		&i.UserID,
	)
	return i, err
}

const getOrderByID = `-- name: GetOrderByID :one
SELECT id, displayed_id, restaurant_id, total_price, status, user_id FROM orders
WHERE id = $1 AND user_id = $2
`

type GetOrderByIDParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetOrderByID(ctx context.Context, arg GetOrderByIDParams) (Order, error) {
	row := q.db.QueryRow(ctx, getOrderByID, arg.ID, arg.UserID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.DisplayedID,
		&i.RestaurantID,
		&i.TotalPrice,
		&i.Status,
		&i.UserID,
	)
	return i, err
}

const getRefresh = `-- name: GetRefresh :one
SELECT crypted_refresh FROM users WHERE id = $1
`

func (q *Queries) GetRefresh(ctx context.Context, id uuid.UUID) (*string, error) {
	row := q.db.QueryRow(ctx, getRefresh, id)
	var crypted_refresh *string
	err := row.Scan(&crypted_refresh)
	return crypted_refresh, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT name, id, phone_number, mail, birthday, created_at FROM users WHERE id = $1
`

type GetUserByIDRow struct {
	Name        string           `json:"name"`
	ID          uuid.UUID        `json:"id"`
	PhoneNumber string           `json:"phone_number"`
	Mail        string           `json:"mail"`
	Birthday    pgtype.Date      `json:"birthday"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (GetUserByIDRow, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i GetUserByIDRow
	err := row.Scan(
		&i.Name,
		&i.ID,
		&i.PhoneNumber,
		&i.Mail,
		&i.Birthday,
		&i.CreatedAt,
	)
	return i, err
}

const getUserOrders = `-- name: GetUserOrders :many
SELECT id, displayed_id, restaurant_id, total_price, status, user_id FROM orders
WHERE user_id = $1
ORDER BY displayed_id
`

func (q *Queries) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	rows, err := q.db.Query(ctx, getUserOrders, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.DisplayedID,
			&i.RestaurantID,
			&i.TotalPrice,
			&i.Status,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const logIn = `-- name: LogIn :one
SELECT id, phone_number, crypted_password, name, mail, birthday, crypted_refresh, created_at FROM users
WHERE phone_number = $1 AND crypted_password = $2
`

type LogInParams struct {
	PhoneNumber     string `json:"phone_number"`
	CryptedPassword string `json:"crypted_password"`
}

func (q *Queries) LogIn(ctx context.Context, arg LogInParams) (User, error) {
	row := q.db.QueryRow(ctx, logIn, arg.PhoneNumber, arg.CryptedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.PhoneNumber,
		&i.CryptedPassword,
		&i.Name,
		&i.Mail,
		&i.Birthday,
		&i.CryptedRefresh,
		&i.CreatedAt,
	)
	return i, err
}

const updateRefresh = `-- name: UpdateRefresh :exec
UPDATE users SET crypted_refresh = $2 WHERE id = $1
`

type UpdateRefreshParams struct {
	ID             uuid.UUID `json:"id"`
	CryptedRefresh *string   `json:"crypted_refresh"`
}

func (q *Queries) UpdateRefresh(ctx context.Context, arg UpdateRefreshParams) error {
	_, err := q.db.Exec(ctx, updateRefresh, arg.ID, arg.CryptedRefresh)
	return err
}
