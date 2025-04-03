package service

import "errors"

var (
	ErrWrongData       = errors.New("incorrect login or password")
	ErrServerError     = errors.New("server error occured")
	ErrWrongToken      = errors.New("wrong token type provided")
	ErrExpiredToken    = errors.New("previous token provided")
	ErrWrongTimeFormat = errors.New("wrong time format provided")
	ErrLoginOccupied   = errors.New("login already occupied")
)

type Tokens struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}
