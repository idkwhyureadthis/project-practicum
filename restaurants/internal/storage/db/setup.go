package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"log/slog"
	"time"

	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/storage/db/generated"
	"github.com/jackc/pgx/v5"
)

func SetupConnection(connUrl, adminPass string) *generated.Queries {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	conn, err := pgx.Connect(ctx, connUrl)
	if err != nil || conn == nil {
		log.Fatal(err)
	}
	sqlcConn := generated.New(conn)
	id, err := sqlcConn.CheckAdmin(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if id == 0 {
		h := sha256.New()
		h.Write([]byte(adminPass))
		err = sqlcConn.SetupAdmin(context.Background(), hex.EncodeToString(h.Sum(nil)))
		if err != nil {
			log.Fatal(err)
		}
		slog.Info("admin created")
	} else {
		slog.Info("admin already created")
	}
	return sqlcConn
}
