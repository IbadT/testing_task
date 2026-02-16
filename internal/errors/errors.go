package myerrors

import (
	"errors"
	domain "testingtask/internal/domain/subscription"
	"testingtask/internal/web/subscriptions"
)

type ErrorResponse struct {
	Error string `json:"error" example:"invalid id format"`
}

type ErrorNotFound struct {
	Error string `json:"error" example:"subscription not found"`
}

type ErrorInternalServerError struct {
	Error string `json:"error" example:"internal server error"`
}

var (
	ErrInvalidID   = errors.New("invalid id format")
	ErrNotFound    = errors.New("subscription not found")
	ErrInvalidData = errors.New("invalid input data")
	ErrInternal    = errors.New("internal server error")

	ErrConflict       = errors.New("conflict")
	ErrDatabase       = errors.New("database error")
	ErrContextTimeout = errors.New("context timeout")
	ErrCreateFailed   = errors.New("failed to create subscription")
	ErrListFailed     = errors.New("failed to list subscriptions")
	ErrUpdateFailed   = errors.New("failed to update subscription")
	ErrDeleteFailed   = errors.New("failed to delete subscription")
)

func MapError(err error) (subscriptions.ErrorResponse, int) {
	switch {
	case errors.Is(err, ErrInvalidID):
		return subscriptions.ErrorResponse{Error: err.Error()}, 400

	case errors.Is(err, ErrInvalidData):
		return subscriptions.ErrorResponse{Error: err.Error()}, 400

	case errors.Is(err, ErrNotFound):
		return subscriptions.ErrorResponse{Error: err.Error()}, 404

	// ДОМЕННЫЕ ОШИБКИ
	case errors.Is(err, domain.ErrInvalidPrice),
		errors.Is(err, domain.ErrInvalidStartDate),
		errors.Is(err, domain.ErrInvalidEndDate),
		errors.Is(err, domain.ErrEmptyServiceName),
		errors.Is(err, domain.ErrCompareDate),
		errors.Is(err, domain.ErrInvalidDate):
		return subscriptions.ErrorResponse{Error: err.Error()}, 400

	// ОШИБКИ РЕПОЗИТОРИЯ
	case errors.Is(err, ErrConflict),
		errors.Is(err, ErrDatabase),
		errors.Is(err, ErrContextTimeout),
		errors.Is(err, ErrCreateFailed),
		errors.Is(err, ErrListFailed),
		errors.Is(err, ErrUpdateFailed),
		errors.Is(err, ErrDeleteFailed):
		return subscriptions.ErrorResponse{Error: err.Error()}, 500

	default:
		return subscriptions.ErrorResponse{Error: "internal server error"}, 500
	}
}
