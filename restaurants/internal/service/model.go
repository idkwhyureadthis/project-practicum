package service

import "errors"

var (
	ErrWrongData    = errors.New("incorrect login or password")
	ErrServerError  = errors.New("server error occured")
	ErrWrongToken   = errors.New("wrong token type provided")
	ErrExpiredToken = errors.New("previous token provided")
)

type Tokens struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}
