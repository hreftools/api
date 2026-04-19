package resource

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	// validation title
	ErrValidationTitleLength            = errors.New("title must be between 3 and 255 characters")
	ErrValidationTitleInvalidCharacters = errors.New("title must not contain control characters")

	// validation description
	ErrValidationDescriptionLength            = errors.New("description must be less than 512 characters")
	ErrValidationDescriptionInvalidCharacters = errors.New("description must not contain control characters")

	// validation url
	ErrValidationURLFormat  = errors.New("url is invalid")
	ErrValidationURLTooLong = errors.New("url must be at most 2048 characters")
	ErrValidationURLPrivate = errors.New("url must not point to a private or local address")

	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type CreateParams struct {
	UserID      uuid.UUID
	Title       string
	Description string
	Url         string
}

type UpdateParams struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Url         string
}

type Repository interface {
	List(ctx context.Context, userID uuid.UUID) ([]Resource, error)
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Resource, error)
	Create(ctx context.Context, params CreateParams) (Resource, error)
	Update(ctx context.Context, params UpdateParams) (Resource, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Resource, error)
}

type Service struct {
	ResourceRepo Repository
}

func NewService(repo Repository) *Service {
	return &Service{ResourceRepo: repo}
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]Resource, error) {
	return s.ResourceRepo.List(ctx, userID)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Resource, error) {
	return s.ResourceRepo.Get(ctx, id, userID)
}

type CreateParamsService struct {
	UserID      uuid.UUID
	Title       string
	Description string
	Url         string
}

func (s *Service) Create(ctx context.Context, params CreateParamsService) (Resource, error) {
	title, err := validateTitle(params.Title)
	if err != nil {
		return Resource{}, err
	}
	description, err := validateDescription(params.Description)
	if err != nil {
		return Resource{}, err
	}
	url, err := validateURL(params.Url)
	if err != nil {
		return Resource{}, err
	}

	repoParams := CreateParams{
		UserID:      params.UserID,
		Title:       title,
		Description: description,
		Url:         url,
	}
	return s.ResourceRepo.Create(ctx, repoParams)
}

type UpdateParamsService struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Url         string
}

func (s *Service) Update(ctx context.Context, params UpdateParamsService) (Resource, error) {
	title, err := validateTitle(params.Title)
	if err != nil {
		return Resource{}, err
	}
	description, err := validateDescription(params.Description)
	if err != nil {
		return Resource{}, err
	}
	url, err := validateURL(params.Url)
	if err != nil {
		return Resource{}, err
	}

	repoParams := UpdateParams{
		ID:          params.ID,
		UserID:      params.UserID,
		Title:       title,
		Description: description,
		Url:         url,
	}
	return s.ResourceRepo.Update(ctx, repoParams)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Resource, error) {
	return s.ResourceRepo.Delete(ctx, id, userID)
}
