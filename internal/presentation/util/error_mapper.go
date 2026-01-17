package util

import (
	"errors"
	"net/http"

	"github.com/yokitheyo/wb_level3_02/internal/application/usecase"
)

// MapErrorToStatus converts usecase errors to HTTP status codes
func MapErrorToStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch {
	case errors.Is(err, usecase.ErrShortCodeAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, usecase.ErrInvalidCustomShort),
		errors.Is(err, usecase.ErrURLRequired),
		errors.Is(err, usecase.ErrInvalidQuery):
		return http.StatusBadRequest
	case errors.Is(err, usecase.ErrNotFound),
		errors.Is(err, usecase.ErrShortCodeRequired):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
