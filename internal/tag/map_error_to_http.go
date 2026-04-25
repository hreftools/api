package tag

import (
	"context"
	"errors"
	"log"
	"net/http"
)

func MapErrorToHTTP(err error) (int, string) {
	// context errors
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusRequestTimeout, "request timeout"
	}
	if errors.Is(err, context.Canceled) {
		return 499, "request cancelled"
	}

	// validation errors
	if errors.Is(err, ErrValidationNameLength) ||
		errors.Is(err, ErrValidationNameCharacters) ||
		errors.Is(err, ErrValidationNameHyphens) ||
		errors.Is(err, ErrValidationTooManyTags) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, "not found"
	}
	if errors.Is(err, ErrConflict) {
		return http.StatusConflict, "conflict"
	}

	log.Printf("Service error: %v", err)
	return http.StatusInternalServerError, "internal server error"
}
