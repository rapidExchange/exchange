package storage

import "errors"

var (
	ErrStockNotFound = errors.New("user not found")
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists = errors.New("user already exists")
)