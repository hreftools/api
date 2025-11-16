package store

import (
	"database/sql"

	"github.com/zapi-sh/api/internal/db"
)

type Store struct {
	Resources ResourceStore
}

func NewStore(pool *sql.DB) *Store {
	queries := db.New(pool)

	return &Store{
		Resources: NewResourceStore(queries),
	}
}
