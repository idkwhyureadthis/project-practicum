package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type VerificationResponse struct {
	Type string `json:"type"`
	Role string `json:"role"`
	Id   string `json:"id"`
}

type TokenClaims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func Generate(role, id string, secret []byte) (*Tokens, error) {
	accessClaims := TokenClaims{
		"access",
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Audience:  []string{role},
			Subject:   id,
		},
	}

	refreshClaims := TokenClaims{
		"refresh",
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Audience:  []string{role},
			Subject:   id,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(secret)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		Access:  accessString,
		Refresh: refreshString,
	}, nil
}

func Verify(token string, secret []byte) (*VerificationResponse, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := parsedToken.Claims.(*TokenClaims); ok && len(claims.Audience) == 1 {
		return &VerificationResponse{
			Id:   claims.Subject,
			Role: claims.Audience[0],
			Type: claims.Type,
		}, nil
	}
	return nil, errors.New("wrong token provided")
}
