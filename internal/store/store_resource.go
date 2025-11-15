package store

import (
	"context"

	"github.com/zapi-sh/api/internal/db"
)

type ResourceStore interface {
	List(ctx context.Context) ([]db.Resource, error)
	Create(ctx context.Context, title string, url string) (db.Resource, error)
	Get(ctx context.Context, id int64) (db.Resource, error)
	Delete(ctx context.Context, id int64) (db.Resource, error)
}

type resourceStore struct {
	queries db.Querier
}

func NewResourceStore(queries db.Querier) ResourceStore {
	return &resourceStore{
		queries: queries,
	}
}

func (r *resourceStore) List(ctx context.Context) ([]db.Resource, error) {
	return r.queries.ListResources(ctx)
}

func (r *resourceStore) Create(ctx context.Context, title string, url string) (db.Resource, error) {
	args := db.CreateResourceParams{
		Title: title,
		Url:   url,
	}
	return r.queries.CreateResource(ctx, args)
}

func (r *resourceStore) Get(ctx context.Context, id int64) (db.Resource, error) {
	return r.queries.GetResource(ctx, id)
}

func (r *resourceStore) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteResource(ctx, id)
}
