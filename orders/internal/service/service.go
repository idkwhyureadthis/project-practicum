package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/idkwhyureadthis/project-practicum/orders/internal/storage/db"
	"github.com/idkwhyureadthis/project-practicum/orders/internal/storage/db/generated"
	"github.com/idkwhyureadthis/project-practicum/orders/pkg/orders"
	"github.com/idkwhyureadthis/project-practicum/orders/pkg/tokens"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db     *pgx.Conn
	conn   *generated.Queries
	secret []byte
}

func New(connUrl string, secretKey string) *Service {
	service := Service{}
	service.conn, service.db = db.SetupConnection(connUrl)
	service.secret = []byte(secretKey)
	return &service
}

func (s *Service) LogIn(phoneNumber, password string) (*tokens.Tokens, *generated.User, error) {
	h := sha256.New()
	h.Write([]byte(password))
	cryptedPass := hex.EncodeToString(h.Sum(nil))

	fmt.Println(phoneNumber, password)
	user, err := s.conn.LogIn(context.Background(), generated.LogInParams{
		PhoneNumber:     phoneNumber,
		CryptedPassword: cryptedPass,
	})
	if err == pgx.ErrNoRows {
		return nil, nil, ErrWrongLoginOrPass
	} else if err != nil {
		return nil, nil, err
	}

	tokens, err := tokens.Generate("user", user.ID.String(), s.secret)
	if err != nil {
		return nil, nil, err
	}

	h = sha256.New()
	h.Write([]byte(tokens.Refresh))
	cryptedRefresh, err := bcrypt.GenerateFromPassword(h.Sum(nil), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	refreshString := string(cryptedRefresh)
	err = s.conn.UpdateRefresh(context.Background(), generated.UpdateRefreshParams{
		ID:             user.ID,
		CryptedRefresh: &refreshString,
	})
	if err != nil {
		return nil, nil, err
	}

	return tokens, &user, nil
}

func (s *Service) SignUp(phoneNumber, password, name, mail string) (*generated.User, *tokens.Tokens, int, error) {
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
		return nil, nil, http.StatusConflict, ErrPhoneOccupied
	}
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	tokensData, err := tokens.Generate("user", user.ID.String(), s.secret)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	h = sha256.New()
	h.Write([]byte(tokensData.Refresh))
	cryptedRefresh, err := bcrypt.GenerateFromPassword(h.Sum(nil), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	refreshString := string(cryptedRefresh)
	err = s.conn.UpdateRefresh(context.Background(), generated.UpdateRefreshParams{
		ID:             user.ID,
		CryptedRefresh: &refreshString,
	})
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	return &user, tokensData, http.StatusOK, nil
}

func (s *Service) Verify(token, expectedType string) (uuid.UUID, error) {
	tokenInfo, err := tokens.Verify(token, s.secret)
	if err != nil {
		return uuid.Nil, ErrWrongToken
	}
	if tokenInfo.Type != expectedType || tokenInfo.Role != "user" {
		return uuid.Nil, ErrWrongToken
	}

	userID, err := uuid.Parse(tokenInfo.Id)
	if err != nil {
		return uuid.Nil, ErrWrongToken
	}

	return userID, nil
}

func (s *Service) Refresh(refreshToken string) (*tokens.Tokens, error) {
	tokenInfo, err := tokens.Verify(refreshToken, s.secret)
	if err != nil {
		return nil, ErrWrongToken
	}

	if tokenInfo.Type != "refresh" || tokenInfo.Role != "user" {
		return nil, ErrWrongToken
	}

	userID, err := uuid.Parse(tokenInfo.Id)
	if err != nil {
		return nil, ErrWrongToken
	}

	storedRefresh, err := s.conn.GetRefresh(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	if storedRefresh == nil {
		return nil, ErrExpiredToken
	}

	h := sha256.New()
	h.Write([]byte(refreshToken))
	if err = bcrypt.CompareHashAndPassword([]byte(*storedRefresh), h.Sum(nil)); err != nil {
		return nil, ErrExpiredToken
	}

	newTokens, err := tokens.Generate("user", tokenInfo.Id, s.secret)
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
	err = s.conn.UpdateRefresh(context.Background(), generated.UpdateRefreshParams{
		ID:             userID,
		CryptedRefresh: &refreshString,
	})
	if err != nil {
		return nil, err
	}

	return newTokens, nil
}

func (s *Service) GetUserByID(id uuid.UUID) (*generated.GetUserByIDRow, error) {
	user, err := s.conn.GetUserByID(context.Background(), id)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, err
}

func (s *Service) InvalidateRefreshToken(id uuid.UUID) error {
	var emptyString string
	return s.conn.UpdateRefresh(context.Background(), generated.UpdateRefreshParams{
		ID:             id,
		CryptedRefresh: &emptyString,
	})
}

func (s *Service) CreateOrder(restaurantID, userID uuid.UUID, itemsID []string, sizes []string) (*generated.Order, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	var err error
	var displayedID int32
	var totalPrice float64
	for !errors.Is(err, pgx.ErrNoRows) {
		displayedID = orders.GenerateOrderId()
		_, err = s.conn.GetOrderByDisplayedId(ctx, generated.GetOrderByDisplayedIdParams{
			RestaurantID: restaurantID,
			DisplayedID:  displayedID,
		})
	}
	defer cancelFunc()
	for idx, itemID := range itemsID {
		itemUUID, err := uuid.Parse(itemID)
		if err != nil {
			return nil, err
		}
		item, err := s.conn.GetItemById(ctx, itemUUID)
		if err != nil {
			return nil, err
		}

		var priceIndex int
		if len(item.Prices) != 1 {
			priceIndex = slices.Index(item.Sizes, sizes[idx])
			if priceIndex == -1 {
				return nil, ErrWrongSize
			}
		}

		totalPrice += item.Prices[priceIndex]
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.conn.WithTx(tx)

	orderId := uuid.New()
	orderData, err := qtx.CreateOrder(ctx, generated.CreateOrderParams{
		ID:           orderId,
		DisplayedID:  displayedID,
		RestaurantID: restaurantID,
		TotalPrice:   totalPrice,
		Status:       "created",
		UserID:       userID,
	})

	if err != nil {
		return nil, err
	}

	for _, id := range itemsID {
		orderUUID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		err = qtx.CreateOrderItem(ctx, generated.CreateOrderItemParams{
			OrderID: orderId,
			ItemID:  orderUUID,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &orderData, nil
}

func (s *Service) GetOrderByID(id, userID uuid.UUID) (*generated.Order, error) {
	order, err := s.conn.GetOrderByID(context.Background(), generated.GetOrderByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *Service) GetUserOrders(userID uuid.UUID) ([]generated.Order, error) {
	return s.conn.GetUserOrders(context.Background(), userID)
}

func (s *Service) DeleteOrder(id, userID uuid.UUID) error {
	return s.conn.DeleteOrder(context.Background(), generated.DeleteOrderParams{
		ID:     id,
		UserID: userID,
	})
}
