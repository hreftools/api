package uow

import (
	"context"

	"github.com/google/uuid"
	"github.com/urlspace/api/internal/collection"
	"github.com/urlspace/api/internal/resource"
	"github.com/urlspace/api/internal/tag"
)

// Repositories groups all repositories available within a transaction. It lives here
// because the uow package coordinates across multiple domain repositories.
// Neither resource nor tag imports this package.
type Repositories struct {
	Resources   resource.Repository
	Tags        tag.Repository
	Collections collection.Repository
}

// UnitOfWork runs fn inside a single database transaction. Every repository in
// the Repositories value passed to fn executes against that transaction.
type UnitOfWork interface {
	RunInTx(ctx context.Context, fn func(Repositories) error) error
}

type Service struct {
	// repos holds repositories for single-repo operations that don't need
	// transactional guarantees. More repositories can be added here as needed.
	repos Repositories
	// uow wraps multi-repo operations in a transaction. Currently used for
	// resource + tag coordination only.
	uow UnitOfWork
}

func NewService(repos Repositories, uow UnitOfWork) *Service {
	return &Service{repos: repos, uow: uow}
}

// CollectionInfo is a lightweight summary of a collection, included in
// enriched resource responses. Only ID and Title are needed for display.
type CollectionInfo struct {
	ID    uuid.UUID
	Title string
}

// EnrichedResource extends resource.Resource with tag and collection data.
// The resource package stays independent of tags and collections, so this
// combined type lives here where all domains are coordinated.
type EnrichedResource struct {
	resource.Resource
	Tags       []string
	Collection *CollectionInfo
}

// collectionInfoFromResource builds a CollectionInfo from the resource's
// JOIN-populated fields (used on Get/List paths).
func collectionInfoFromResource(r resource.Resource) *CollectionInfo {
	if r.CollectionID == nil {
		return nil
	}
	return &CollectionInfo{ID: *r.CollectionID, Title: r.CollectionTitle}
}

type CreateResourceParams struct {
	UserID       uuid.UUID
	Title        string
	Description  string
	URL          string
	CollectionID *uuid.UUID
	Tags         []string
}

func (s *Service) CreateResource(ctx context.Context, params CreateResourceParams) (EnrichedResource, error) {
	title, err := resource.ValidateTitle(params.Title)
	if err != nil {
		return EnrichedResource{}, err
	}
	description, err := resource.ValidateDescription(params.Description)
	if err != nil {
		return EnrichedResource{}, err
	}
	url, err := resource.ValidateURL(params.URL)
	if err != nil {
		return EnrichedResource{}, err
	}
	tagNames, err := tag.ValidateTagNames(params.Tags)
	if err != nil {
		return EnrichedResource{}, err
	}

	var result EnrichedResource

	err = s.uow.RunInTx(ctx, func(repos Repositories) error {
		// Validate collection ownership if provided.
		if params.CollectionID != nil {
			c, err := repos.Collections.Get(ctx, *params.CollectionID, params.UserID)
			if err != nil {
				return err
			}
			result.Collection = &CollectionInfo{ID: c.ID, Title: c.Title}
		}

		r, err := repos.Resources.Create(ctx, resource.CreateParams{
			UserID:       params.UserID,
			Title:        title,
			Description:  description,
			URL:          url,
			CollectionID: params.CollectionID,
		})
		if err != nil {
			return err
		}
		result.Resource = r

		tagIDs := make([]uuid.UUID, 0, len(tagNames))
		for _, name := range tagNames {
			t, err := repos.Tags.UpsertByName(ctx, params.UserID, name)
			if err != nil {
				return err
			}
			tagIDs = append(tagIDs, t.ID)
		}

		if err := repos.Tags.ReplaceResourceTags(ctx, r.ID, tagIDs); err != nil {
			return err
		}

		tags, err := repos.Tags.GetTagsForResource(ctx, r.ID)
		if err != nil {
			return err
		}
		result.Tags = tags

		return nil
	})

	return result, err
}

type UpdateResourceParams struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Title        string
	Description  string
	URL          string
	CollectionID *uuid.UUID
	Tags         []string
}

func (s *Service) UpdateResource(ctx context.Context, params UpdateResourceParams) (EnrichedResource, error) {
	title, err := resource.ValidateTitle(params.Title)
	if err != nil {
		return EnrichedResource{}, err
	}
	description, err := resource.ValidateDescription(params.Description)
	if err != nil {
		return EnrichedResource{}, err
	}
	url, err := resource.ValidateURL(params.URL)
	if err != nil {
		return EnrichedResource{}, err
	}
	tagNames, err := tag.ValidateTagNames(params.Tags)
	if err != nil {
		return EnrichedResource{}, err
	}

	var result EnrichedResource

	err = s.uow.RunInTx(ctx, func(repos Repositories) error {
		// Validate collection ownership if provided.
		if params.CollectionID != nil {
			c, err := repos.Collections.Get(ctx, *params.CollectionID, params.UserID)
			if err != nil {
				return err
			}
			result.Collection = &CollectionInfo{ID: c.ID, Title: c.Title}
		}

		r, err := repos.Resources.Update(ctx, resource.UpdateParams{
			ID:           params.ID,
			UserID:       params.UserID,
			Title:        title,
			Description:  description,
			URL:          url,
			CollectionID: params.CollectionID,
		})
		if err != nil {
			return err
		}
		result.Resource = r

		tagIDs := make([]uuid.UUID, 0, len(tagNames))
		for _, name := range tagNames {
			t, err := repos.Tags.UpsertByName(ctx, params.UserID, name)
			if err != nil {
				return err
			}
			tagIDs = append(tagIDs, t.ID)
		}

		if err := repos.Tags.ReplaceResourceTags(ctx, r.ID, tagIDs); err != nil {
			return err
		}

		tags, err := repos.Tags.GetTagsForResource(ctx, r.ID)
		if err != nil {
			return err
		}
		result.Tags = tags

		return nil
	})

	return result, err
}

func (s *Service) ListResources(ctx context.Context, userID uuid.UUID) ([]EnrichedResource, error) {
	list, err := s.repos.Resources.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return []EnrichedResource{}, nil
	}

	resourceIDs := make([]uuid.UUID, len(list))
	for i, item := range list {
		resourceIDs[i] = item.ID
	}

	tagsMap, err := s.repos.Tags.GetTagsForResources(ctx, resourceIDs)
	if err != nil {
		return nil, err
	}

	result := make([]EnrichedResource, len(list))
	for i, item := range list {
		tags := tagsMap[item.ID]
		if tags == nil {
			tags = []string{}
		}
		result[i] = EnrichedResource{
			Resource:   item,
			Tags:       tags,
			Collection: collectionInfoFromResource(item),
		}
	}

	return result, nil
}

func (s *Service) GetResource(ctx context.Context, id uuid.UUID, userID uuid.UUID) (EnrichedResource, error) {
	r, err := s.repos.Resources.Get(ctx, id, userID)
	if err != nil {
		return EnrichedResource{}, err
	}

	tags, err := s.repos.Tags.GetTagsForResource(ctx, id)
	if err != nil {
		return EnrichedResource{}, err
	}

	return EnrichedResource{
		Resource:   r,
		Tags:       tags,
		Collection: collectionInfoFromResource(r),
	}, nil
}

func (s *Service) DeleteResource(ctx context.Context, id uuid.UUID, userID uuid.UUID) (EnrichedResource, error) {
	tags, err := s.repos.Tags.GetTagsForResource(ctx, id)
	if err != nil {
		return EnrichedResource{}, err
	}

	// Look up collection info before deleting (DELETE can't JOIN).
	r, err := s.repos.Resources.Get(ctx, id, userID)
	if err != nil {
		return EnrichedResource{}, err
	}
	col := collectionInfoFromResource(r)

	deleted, err := s.repos.Resources.Delete(ctx, id, userID)
	if err != nil {
		return EnrichedResource{}, err
	}

	return EnrichedResource{
		Resource:   deleted,
		Tags:       tags,
		Collection: col,
	}, nil
}
