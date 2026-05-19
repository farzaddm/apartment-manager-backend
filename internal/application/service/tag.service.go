package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"

	"github.com/google/uuid"
)

type TagService interface {
	Create(ctx context.Context, req dto.CreateTagRequest) (*dto.TagResponse, error)
	GetAll(ctx context.Context) ([]dto.TagResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagService struct {
	tagRepo postgres.TagInterface
}

func NewTagService(tagRepo postgres.TagInterface) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) Create(ctx context.Context, req dto.CreateTagRequest) (*dto.TagResponse, error) {
	tag := &entity.Tag{Name: req.Name}
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return &dto.TagResponse{ID: tag.ID, Name: tag.Name}, nil
}

func (s *tagService) GetAll(ctx context.Context) ([]dto.TagResponse, error) {
	tags, err := s.tagRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TagResponse, len(tags))
	for i, t := range tags {
		res[i] = dto.TagResponse{ID: t.ID, Name: t.Name}
	}
	return res, nil
}

func (s *tagService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.tagRepo.Delete(ctx, id)
}
