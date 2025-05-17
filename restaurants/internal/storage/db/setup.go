package db

import (
	"context"
	"log"
	"time"

	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/storage/db/generated"
	"github.com/jackc/pgx/v5"
)

func SetupConnection(connUrl string) *generated.Queries {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, connUrl)
	if err != nil || conn == nil {
		log.Fatal(err)
	}
	sqlcConn := generated.New(conn)
	return sqlcConn
}
