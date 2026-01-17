package usecase

import "errors"

var (
	ErrShortCodeAlreadyExists = errors.New("custom short code already exists")
	ErrInvalidCustomShort     = errors.New("invalid custom short code")
	ErrURLRequired            = errors.New("URL is required")
	ErrShortCodeRequired      = errors.New("short code is required")
	ErrNotFound               = errors.New("not found")
	ErrInvalidQuery           = errors.New("invalid query parameters")
)
