package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/storage/db"
	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/storage/db/generated"
	"github.com/idkwhyureadthis/project-practicum/restaurants/pkg/tokens"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	conn   *generated.Queries
	secret []byte
}

func New(connUrl, adminPass, secret string) *Service {
	service := Service{}
	service.conn = db.SetupConnection(connUrl, adminPass)
	service.secret = []byte(secret)
	return &service
}

func (s *Service) LogIn(login, password string) (*tokens.Tokens, error) {
	h := sha256.New()
	h.Write([]byte(password))
	cryptedPass := hex.EncodeToString(h.Sum(nil))
	data, err := s.conn.GetAdmin(context.Background(), generated.GetAdminParams{
		Login:           login,
		CryptedPassword: cryptedPass,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrWrongData
	} else if err != nil {
		return nil, ErrServerError
	}

	role := "admin"
	if data.IsSuperadmin {
		role = "superadmin"
	}
	tokens, err := tokens.Generate(role, data.ID.String(), s.secret)
	if err != nil {
		return nil, err
	}
	h = sha256.New()
	h.Write([]byte(tokens.Refresh))
	cryptedRefresh, err := bcrypt.GenerateFromPassword(h.Sum(nil), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	refreshString := string(cryptedRefresh)
	s.conn.UpdateRefresh(context.Background(), generated.UpdateRefreshParams{
		CryptedRefresh: &refreshString,
		ID:             data.ID,
	})

	return tokens, nil
}

func (s *Service) Verify(token, expectedRole string) (*string, error) {
	tokenInfo, err := tokens.Verify(token, s.secret)
	if err != nil {
		return nil, err
	}
	if tokenInfo.Type != expectedRole {
		return nil, ErrWrongToken
	}
	return &tokenInfo.Role, nil
}

func (s *Service) Generate(refresh string) (*tokens.Tokens, error) {
	tokenInfo, err := tokens.Verify(refresh, s.secret)
	if err != nil {
		return nil, err
	}

	if tokenInfo.Type != "refresh" {
		return nil, ErrWrongToken
	}

	storedRefresh, err := s.conn.GetRefresh(context.Background(), uuid.MustParse(tokenInfo.Id))

	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write([]byte(refresh))
	if err = bcrypt.CompareHashAndPassword([]byte(*storedRefresh), h.Sum(nil)); err != nil {
		return nil, ErrExpiredToken
	}

	newTokens, err := tokens.Generate(tokenInfo.Role, tokenInfo.Id, s.secret)
	if err != nil {
		return nil, err
	}

	h = sha256.New()
	h.Write([]byte(newTokens.Refresh))
	cryptedRefresh, err := bcrypt.GenerateFromPassword(h.Sum(nil), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	refreshString := string(cryptedRefresh)
	s.conn.UpdateRefresh(context.Background(), generated.UpdateRefreshParams{
		CryptedRefresh: &refreshString,
		ID:             uuid.MustParse(tokenInfo.Id),
	})

	return newTokens, nil
}
