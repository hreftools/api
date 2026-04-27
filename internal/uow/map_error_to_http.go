package uow

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/urlspace/api/internal/collection"
	"github.com/urlspace/api/internal/resource"
	"github.com/urlspace/api/internal/tag"
)

// MapErrorToHTTP maps errors from the uow service to HTTP status codes.
// This covers both resource and tag validation errors since the uow service
// coordinates across both domains.
func MapErrorToHTTP(err error) (int, string) {
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusRequestTimeout, "request timeout"
	}
	if errors.Is(err, context.Canceled) {
		return 499, "request cancelled"
	}

	// resource validation errors
	if errors.Is(err, resource.ErrValidationTitleLength) ||
		errors.Is(err, resource.ErrValidationTitleInvalidCharacters) ||
		errors.Is(err, resource.ErrValidationDescriptionLength) ||
		errors.Is(err, resource.ErrValidationDescriptionInvalidCharacters) ||
		errors.Is(err, resource.ErrValidationURLFormat) ||
		errors.Is(err, resource.ErrValidationURLTooLong) ||
		errors.Is(err, resource.ErrValidationURLPrivate) {
		return http.StatusBadRequest, err.Error()
	}

	// tag validation errors
	if errors.Is(err, tag.ErrValidationNameLength) ||
		errors.Is(err, tag.ErrValidationNameCharacters) ||
		errors.Is(err, tag.ErrValidationNameHyphens) ||
		errors.Is(err, tag.ErrValidationTooManyTags) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, resource.ErrNotFound) || errors.Is(err, tag.ErrNotFound) || errors.Is(err, collection.ErrNotFound) {
		return http.StatusNotFound, "not found"
	}
	if errors.Is(err, resource.ErrConflict) || errors.Is(err, tag.ErrConflict) || errors.Is(err, collection.ErrConflict) {
		return http.StatusConflict, "conflict"
	}

	log.Printf("Service error: %v", err)
	return http.StatusInternalServerError, "internal server error"
}
