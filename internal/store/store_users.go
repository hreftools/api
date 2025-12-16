package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/zapi-sh/api/internal/db"
)

const (
	UserUsernameLengthMin = 3
	UserUsernameLengthMax = 32
	UserPasswordLengthMin = 12
)

type UserStore interface {
	List(ctx context.Context) ([]db.User, error)
	Get(ctx context.Context, id uuid.UUID) (db.User, error)
	Create(ctx context.Context, username string, email string, password string) (db.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userStore struct {
	queries db.Querier
}

func NewUserStore(queries db.Querier) UserStore {
	return &userStore{
		queries: queries,
	}
}

func (r *userStore) List(ctx context.Context) ([]db.User, error) {
	return r.queries.ListUsers(ctx)
}

func (r *userStore) Get(ctx context.Context, id uuid.UUID) (db.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *userStore) Create(ctx context.Context, username string, email string, password string) (db.User, error) {
	args := db.CreateUserParams{
		Username: username,
		Email:    email,
		Password: password,
	}

	return r.queries.CreateUser(ctx, args)
}

func (r *userStore) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}
