package uow

import (
	"context"

	"github.com/google/uuid"
	"github.com/urlspace/api/internal/resource"
	"github.com/urlspace/api/internal/tag"
)

// Repositories groups all repositories available within a transaction. It lives here
// because the uow package coordinates across multiple domain repositories.
// Neither resource nor tag imports this package.
type Repositories struct {
	Resources resource.Repository
	Tags      tag.Repository
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

// ResourceWithTags extends resource.Resource with tag data. The resource
// package stays independent of tags, so this combined type lives here
// where both domains are coordinated.
type ResourceWithTags struct {
	resource.Resource
	Tags []string
}

type CreateResourceParams struct {
	UserID      uuid.UUID
	Title       string
	Description string
	Url         string
	Tags        []string
}

func (s *Service) CreateResource(ctx context.Context, params CreateResourceParams) (ResourceWithTags, error) {
	title, err := resource.ValidateTitle(params.Title)
	if err != nil {
		return ResourceWithTags{}, err
	}
	description, err := resource.ValidateDescription(params.Description)
	if err != nil {
		return ResourceWithTags{}, err
	}
	url, err := resource.ValidateURL(params.Url)
	if err != nil {
		return ResourceWithTags{}, err
	}
	tagNames, err := tag.ValidateTagNames(params.Tags)
	if err != nil {
		return ResourceWithTags{}, err
	}

	var result ResourceWithTags

	err = s.uow.RunInTx(ctx, func(repos Repositories) error {
		r, err := repos.Resources.Create(ctx, resource.CreateParams{
			UserID:      params.UserID,
			Title:       title,
			Description: description,
			URL:         url,
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
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Url         string
	Tags        []string
}

func (s *Service) UpdateResource(ctx context.Context, params UpdateResourceParams) (ResourceWithTags, error) {
	title, err := resource.ValidateTitle(params.Title)
	if err != nil {
		return ResourceWithTags{}, err
	}
	description, err := resource.ValidateDescription(params.Description)
	if err != nil {
		return ResourceWithTags{}, err
	}
	url, err := resource.ValidateURL(params.Url)
	if err != nil {
		return ResourceWithTags{}, err
	}
	tagNames, err := tag.ValidateTagNames(params.Tags)
	if err != nil {
		return ResourceWithTags{}, err
	}

	var result ResourceWithTags

	err = s.uow.RunInTx(ctx, func(repos Repositories) error {
		r, err := repos.Resources.Update(ctx, resource.UpdateParams{
			ID:          params.ID,
			UserID:      params.UserID,
			Title:       title,
			Description: description,
			URL:         url,
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

func (s *Service) ListResources(ctx context.Context, userID uuid.UUID) ([]ResourceWithTags, error) {
	list, err := s.repos.Resources.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return []ResourceWithTags{}, nil
	}

	resourceIDs := make([]uuid.UUID, len(list))
	for i, item := range list {
		resourceIDs[i] = item.ID
	}

	tagsMap, err := s.repos.Tags.GetTagsForResources(ctx, resourceIDs)
	if err != nil {
		return nil, err
	}

	result := make([]ResourceWithTags, len(list))
	for i, item := range list {
		tags := tagsMap[item.ID]
		if tags == nil {
			tags = []string{}
		}
		result[i] = ResourceWithTags{Resource: item, Tags: tags}
	}

	return result, nil
}

func (s *Service) GetResource(ctx context.Context, id uuid.UUID, userID uuid.UUID) (ResourceWithTags, error) {
	r, err := s.repos.Resources.Get(ctx, id, userID)
	if err != nil {
		return ResourceWithTags{}, err
	}

	tags, err := s.repos.Tags.GetTagsForResource(ctx, id)
	if err != nil {
		return ResourceWithTags{}, err
	}

	return ResourceWithTags{Resource: r, Tags: tags}, nil
}

func (s *Service) DeleteResource(ctx context.Context, id uuid.UUID, userID uuid.UUID) (ResourceWithTags, error) {
	tags, err := s.repos.Tags.GetTagsForResource(ctx, id)
	if err != nil {
		return ResourceWithTags{}, err
	}

	r, err := s.repos.Resources.Delete(ctx, id, userID)
	if err != nil {
		return ResourceWithTags{}, err
	}

	return ResourceWithTags{Resource: r, Tags: tags}, nil
}
