package db

import (
	"context"
	"log"

	"github.com/idkwhyureadthis/project-practicum/restaurant-service/internal/storage/db/generated"
	"github.com/jackc/pgx/v5"
)

func SetupConnection(connUrl string) *generated.Queries {
	// TODO: define context
	conn, err := pgx.Connect(context.TODO(), connUrl)
	if err != nil || conn == nil {
		log.Fatal(err)
	}
	return generated.New(conn)
}
