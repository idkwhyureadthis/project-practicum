package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/storage/db"
	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/storage/db/generated"
	"github.com/idkwhyureadthis/project-practicum/restaurants/pkg/timeconverter"
	"github.com/idkwhyureadthis/project-practicum/restaurants/pkg/tokens"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *Service) AddRestaurant(openTime, closeTime, name string, latitude, longitude float64) (*uuid.UUID, error) {
	pointCoords := pgtype.Vec2{X: latitude, Y: longitude}
	point := pgtype.Point{P: pointCoords, Valid: true}

	parsedOpenTime, err := time.Parse(time.TimeOnly, openTime)
	if err != nil {
		return nil, ErrWrongTimeFormat
	}

	parsedCloseTime, err := time.Parse(time.TimeOnly, closeTime)
	if err != nil {
		return nil, ErrWrongTimeFormat
	}

	mcsOpenTime := timeconverter.TimeToMicro(parsedOpenTime)
	mcsCloseTime := timeconverter.TimeToMicro(parsedCloseTime)

	pgOpenTime := pgtype.Time{Valid: true, Microseconds: mcsOpenTime}
	pgCloseTime := pgtype.Time{Valid: true, Microseconds: mcsCloseTime}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	restId, err := s.conn.AddRestaurant(ctx, generated.AddRestaurantParams{
		Coordinates: point,
		OpenTime:    pgOpenTime,
		CloseTime:   pgCloseTime,
		Name:        name,
	})
	cancel()
	if err != nil {
		return nil, err
	}
	return &restId, nil
}

func (s *Service) GetRestaurants() (*[]generated.Restaurant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	rests, err := s.conn.GetRestaurants(ctx)
	cancel()
	if err != nil {
		return nil, err
	}
	return &rests, nil
}

func (s *Service) CreateAdmin(restaurantId uuid.UUID, login, password string) (*uuid.UUID, error) {
	h := sha256.New()
	h.Write([]byte(password))
	cryptedPass := hex.EncodeToString(h.Sum(nil))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	adminId, err := s.conn.CreateAdmin(ctx, generated.CreateAdminParams{
		Login:           login,
		CryptedPassword: cryptedPass,
		RestaurantID:    &restaurantId,
	})
	cancel()
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, ErrLoginOccupied
	}
	if err != nil {
		return nil, err
	}

	return &adminId, nil
}

func (s *Service) CreateItem(sizes, prices []string, name, description string, images []*multipart.FileHeader) (*uuid.UUID, error) {
	crypytedImages := []string{}
	for _, file := range images {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		bytes, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}
		ext := filepath.Ext(file.Filename)[1:]
		photoString := base64.StdEncoding.EncodeToString(bytes)
		photoString = fmt.Sprintf("data:image/%s; base64,", ext) + photoString
		crypytedImages = append(crypytedImages, photoString)
	}
	parsedPrices := []float64{}
	for _, elem := range prices {
		parsedFloat, err := strconv.ParseFloat(elem, 64)
		if err != nil {
			return nil, err
		}
		parsedPrices = append(parsedPrices, parsedFloat)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	uuid, err := s.conn.CreateItem(ctx, generated.CreateItemParams{
		Photos:      crypytedImages,
		Name:        name,
		Description: description,
		Prices:      parsedPrices,
		Sizes:       sizes,
	})
	cancel()
	if err != nil {
		return nil, err
	}

	return &uuid, nil
}
