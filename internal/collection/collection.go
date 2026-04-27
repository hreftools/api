package collection

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Public      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	// validation title
	ErrValidationTitleLength            = errors.New("title must be between 3 and 255 characters")
	ErrValidationTitleInvalidCharacters = errors.New("title must not contain control characters")

	// validation description
	ErrValidationDescriptionLength            = errors.New("description must be less than 512 characters")
	ErrValidationDescriptionInvalidCharacters = errors.New("description must not contain control characters")

	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type CreateParams struct {
	UserID      uuid.UUID
	Title       string
	Description string
	Public      bool
}

type UpdateParams struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Public      bool
}

type Repository interface {
	List(ctx context.Context, userID uuid.UUID) ([]Collection, error)
	Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Collection, error)
	Create(ctx context.Context, params CreateParams) (Collection, error)
	Update(ctx context.Context, params UpdateParams) (Collection, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Collection, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]Collection, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Collection, error) {
	return s.repo.Get(ctx, id, userID)
}

func (s *Service) Create(ctx context.Context, params CreateParams) (Collection, error) {
	title, err := ValidateTitle(params.Title)
	if err != nil {
		return Collection{}, err
	}
	description, err := ValidateDescription(params.Description)
	if err != nil {
		return Collection{}, err
	}

	return s.repo.Create(ctx, CreateParams{
		UserID:      params.UserID,
		Title:       title,
		Description: description,
		Public:      params.Public,
	})
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (Collection, error) {
	title, err := ValidateTitle(params.Title)
	if err != nil {
		return Collection{}, err
	}
	description, err := ValidateDescription(params.Description)
	if err != nil {
		return Collection{}, err
	}

	return s.repo.Update(ctx, UpdateParams{
		ID:          params.ID,
		UserID:      params.UserID,
		Title:       title,
		Description: description,
		Public:      params.Public,
	})
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Collection, error) {
	return s.repo.Delete(ctx, id, userID)
}
