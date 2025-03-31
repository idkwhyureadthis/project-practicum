package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/idkwhyureadthis/project-practicum/orders/internal/storage/db"
	"github.com/idkwhyureadthis/project-practicum/orders/internal/storage/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Service struct {
	conn *generated.Queries
}

func New(connUrl string) *Service {
	service := Service{}
	service.conn = db.SetupConnection(connUrl)
	return &service
}

func (s *Service) LogIn(phoneNumber, password string) (*generated.User, int, error) {
	h := sha256.New()
	h.Write([]byte(password))
	crypytedPass := hex.EncodeToString(h.Sum(nil))

	fmt.Println(phoneNumber, password)
	user, err := s.conn.LogIn(context.Background(), generated.LogInParams{
		PhoneNumber:     phoneNumber,
		CryptedPassword: crypytedPass,
	})
	if err == pgx.ErrNoRows {
		return nil, 401, err
	} else if err != nil {
		return nil, 400, err
	}
	return &user, 201, nil
}

func (s *Service) SignUp(phoneNumber, password, name, mail string) (*generated.User, int, error) {
	h := sha256.New()
	h.Write([]byte(password))
	cryptedPass := hex.EncodeToString(h.Sum(nil))

	user, err := s.conn.CreateUser(context.Background(), generated.CreateUserParams{
		PhoneNumber:     phoneNumber,
		CryptedPassword: cryptedPass,
		Name:            name,
		Mail:            mail,
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, 409, ErrPhoneOccupied
	}

	return &user, 200, nil
}
