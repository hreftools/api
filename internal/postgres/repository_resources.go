package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/urlspace/api/internal/db"
	"github.com/urlspace/api/internal/resource"
)

type ResourceRepository struct {
	queries db.Querier
}

func NewResourceRepository(queries db.Querier) resource.Repository {
	return &ResourceRepository{queries: queries}
}

func translateResourceError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return resource.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return resource.ErrConflict
	}
	return err
}

func toCollectionID(n uuid.NullUUID) *uuid.UUID {
	if n.Valid {
		return &n.UUID
	}
	return nil
}

func toNullUUID(id *uuid.UUID) uuid.NullUUID {
	if id != nil {
		return uuid.NullUUID{UUID: *id, Valid: true}
	}
	return uuid.NullUUID{}
}

// toResource maps a db.Resource to a domain Resource. Used by Create, Update,
// and Delete which return plain table columns via RETURNING *. Get and List
// use a custom mapping because their LEFT JOIN returns additional columns
// (CollectionTitle) not present in db.Resource.
func toResource(r db.Resource) resource.Resource {
	return resource.Resource{
		ID:           r.ID,
		UserID:       r.UserID,
		Title:        r.Title,
		Description:  r.Description,
		URL:          r.Url,
		CollectionID: toCollectionID(r.CollectionID),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func (r *ResourceRepository) List(ctx context.Context, userID uuid.UUID) ([]resource.Resource, error) {
	rows, err := r.queries.ListResources(ctx, userID)
	if err != nil {
		return nil, translateResourceError(err)
	}

	resources := make([]resource.Resource, len(rows))
	for i, row := range rows {
		resources[i] = resource.Resource{
			ID:              row.ID,
			UserID:          row.UserID,
			Title:           row.Title,
			Description:     row.Description,
			URL:             row.Url,
			CollectionID:    toCollectionID(row.CollectionID),
			CollectionTitle: row.CollectionTitle.String,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
		}
	}
	return resources, nil
}

func (r *ResourceRepository) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (resource.Resource, error) {
	row, err := r.queries.GetResource(ctx, db.GetResourceParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return resource.Resource{}, translateResourceError(err)
	}
	return resource.Resource{
		ID:              row.ID,
		UserID:          row.UserID,
		Title:           row.Title,
		Description:     row.Description,
		URL:             row.Url,
		CollectionID:    toCollectionID(row.CollectionID),
		CollectionTitle: row.CollectionTitle.String,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}, nil
}

func (r *ResourceRepository) Create(ctx context.Context, params resource.CreateParams) (resource.Resource, error) {
	args := db.CreateResourceParams{
		UserID:       params.UserID,
		Title:        params.Title,
		Description:  params.Description,
		Url:          params.URL,
		CollectionID: toNullUUID(params.CollectionID),
	}
	row, err := r.queries.CreateResource(ctx, args)
	if err != nil {
		return resource.Resource{}, translateResourceError(err)
	}
	return toResource(row), nil
}

func (r *ResourceRepository) Update(ctx context.Context, params resource.UpdateParams) (resource.Resource, error) {
	args := db.UpdateResourceParams{
		ID:           params.ID,
		UserID:       params.UserID,
		Title:        params.Title,
		Description:  params.Description,
		Url:          params.URL,
		CollectionID: toNullUUID(params.CollectionID),
	}
	row, err := r.queries.UpdateResource(ctx, args)
	if err != nil {
		return resource.Resource{}, translateResourceError(err)
	}
	return toResource(row), nil
}

func (r *ResourceRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) (resource.Resource, error) {
	row, err := r.queries.DeleteResource(ctx, db.DeleteResourceParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return resource.Resource{}, translateResourceError(err)
	}
	return toResource(row), nil
}
