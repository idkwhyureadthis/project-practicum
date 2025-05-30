package db

import (
	"context"
	"log"

	"github.com/idkwhyureadthis/project-practicum/orders/internal/storage/db/generated"
	"github.com/jackc/pgx/v5"
)

func SetupConnection(connUrl string) (*generated.Queries, *pgx.Conn) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil || conn == nil {
		log.Fatal(err)
	}
	return generated.New(conn), conn
}
