package store

import (
	"database/sql"

	"github.com/zapi-sh/api/internal/db"
)

type Store struct {
	Resources ResourceStore
}

func NewStore(connectoin *sql.DB) *Store {
	queries := db.New(connectoin)

	return &Store{
		Resources: NewResourceStore(queries),
	}
}
