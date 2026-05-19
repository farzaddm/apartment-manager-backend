package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"

	"github.com/google/uuid"
)

type AnnouncementService interface {
	Create(ctx context.Context, apartmentID uuid.UUID, req dto.CreateAnnouncementRequest) (*dto.AnnouncementResponse, error)
	GetByID(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) (*dto.AnnouncementResponse, error)
	Update(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID, req dto.UpdateAnnouncementRequest) (*dto.AnnouncementResponse, error)
	Delete(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) error
}

type announcementService struct {
	announcementRepo postgres.AnnouncementInterface
	tagRepo          postgres.TagInterface
}

func NewAnnouncementService(ar postgres.AnnouncementInterface, tr postgres.TagInterface) AnnouncementService {
	return &announcementService{announcementRepo: ar, tagRepo: tr}
}

func (s *announcementService) Create(ctx context.Context, apartmentID uuid.UUID, req dto.CreateAnnouncementRequest) (*dto.AnnouncementResponse, error) {
	announcementID := uuid.New() // Pre-generate UUID so we can assign it to our join structures

	var joinTags []entity.TicketAnnouncementTag
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.FindByIDs(ctx, req.TagIDs)
		if err != nil {
			return nil, err
		}

		joinTags = make([]entity.TicketAnnouncementTag, len(tags))
		for i, t := range tags {
			joinTags[i] = entity.TicketAnnouncementTag{
				TagID:          t.ID,
				AnnouncementID: &announcementID,
				Tag:            t, // Attach the full tag model for immediate mapping representation in responses
			}
		}
	}

	announcement := &entity.Announcement{
		BaseModel:   entity.BaseModel{ID: announcementID},
		ApartmentID: apartmentID,
		Title:       req.Title,
		Description: req.Description,
		Body:        req.Body,
		Order:       req.Order,
		IsPinned:    req.IsPinned,
		ExpiredDate: req.ExpiredDate,
		Tags:        joinTags,
	}

	if err := s.announcementRepo.Create(ctx, announcement); err != nil {
		return nil, err
	}

	return s.mapToResponse(announcement), nil
}

func (s *announcementService) GetByID(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) (*dto.AnnouncementResponse, error) {
	announcement, err := s.announcementRepo.FindByIDAndApartment(ctx, id, apartmentID)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(announcement), nil
}

func (s *announcementService) Update(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID, req dto.UpdateAnnouncementRequest) (*dto.AnnouncementResponse, error) {
	announcement, err := s.announcementRepo.FindByIDAndApartment(ctx, id, apartmentID)
	if err != nil {
		return nil, err
	}

	announcement.Title = req.Title
	announcement.Description = req.Description
	announcement.Body = req.Body
	announcement.Order = req.Order
	announcement.IsPinned = req.IsPinned
	announcement.ExpiredDate = req.ExpiredDate

	if err := s.announcementRepo.Update(ctx, announcement); err != nil {
		return nil, err
	}

	var joinTags []entity.TicketAnnouncementTag
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.FindByIDs(ctx, req.TagIDs)
		if err != nil {
			return nil, err
		}

		joinTags = make([]entity.TicketAnnouncementTag, len(tags))
		for i, t := range tags {
			joinTags[i] = entity.TicketAnnouncementTag{
				TagID:          t.ID,
				AnnouncementID: &id,
				Tag:            t,
			}
		}
	}

	if err := s.announcementRepo.ReplaceTags(ctx, id, joinTags); err != nil {
		return nil, err
	}
	announcement.Tags = joinTags

	return s.mapToResponse(announcement), nil
}

func (s *announcementService) Delete(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) error {
	return s.announcementRepo.Delete(ctx, id, apartmentID)
}

func (s *announcementService) mapToResponse(a *entity.Announcement) *dto.AnnouncementResponse {
	tagsRes := make([]dto.TagResponse, len(a.Tags))
	for i, jt := range a.Tags {
		tagsRes[i] = dto.TagResponse{
			ID:   jt.TagID,
			Name: jt.Tag.Name,
		}
	}

	return &dto.AnnouncementResponse{
		ID:          a.ID,
		ApartmentID: a.ApartmentID,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
		Order:       a.Order,
		IsPinned:    a.IsPinned,
		ExpiredDate: a.ExpiredDate,
		Tags:        tagsRes,
		CreatedAt:   a.CreatedAt,
	}
}
