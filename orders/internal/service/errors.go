package service

import "errors"

var ErrNotFound = errors.New("user with such phone number not found")

var ErrWrongLoginOrPass = errors.New("wrong login or password provided")

var ErrPhoneOccupied = errors.New("user with such phone number already exists")
