package service

import (
	"github.com/idkwhyureadthis/project-practicum/orders/internal/storage/db"
	"github.com/idkwhyureadthis/project-practicum/orders/internal/storage/db/generated"
)

type Service struct {
	conn *generated.Queries
}

func New(connUrl string) *Service {
	service := Service{}
	service.conn = db.SetupConnection(connUrl)
	return &service
}
