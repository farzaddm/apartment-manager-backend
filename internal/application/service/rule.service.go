package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"errors"

	"github.com/google/uuid"
)

type RuleService struct {
	ruleRepo postgres.RuleRepository
}

func NewRuleService(ruleRepo postgres.RuleRepository) *RuleService {
	return &RuleService{ruleRepo: ruleRepo}
}

func (s *RuleService) Create(ctx context.Context, apartmentID uuid.UUID, req dto.CreateRuleRequest) (*dto.RuleResponse, error) {
	rule := &entity.Rule{
		BaseModel:   entity.BaseModel{ID: uuid.New()},
		ApartmentID: apartmentID,
		Title:       req.Title,
		Description: req.Description,
		Category:    entity.RuleCategory(req.Category),
	}

	if err := s.ruleRepo.Create(ctx, rule); err != nil {
		return nil, err
	}

	return s.mapToResponse(rule), nil
}

func (s *RuleService) GetByID(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) (*dto.RuleResponse, error) {
	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rule == nil || rule.ApartmentID != apartmentID {
		return nil, errors.New("rule not found")
	}
	return s.mapToResponse(rule), nil
}

func (s *RuleService) ListByApartmentID(ctx context.Context, apartmentID uuid.UUID) ([]dto.RuleResponse, error) {
	rules, err := s.ruleRepo.GetByApartmentID(ctx, apartmentID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.RuleResponse, len(rules))
	for i, r := range rules {
		resp[i] = *s.mapToResponse(&r)
	}
	return resp, nil
}

func (s *RuleService) ListByApartmentAndCategory(ctx context.Context, apartmentID uuid.UUID, category string) ([]dto.RuleResponse, error) {
	rules, err := s.ruleRepo.GetByApartmentAndCategory(ctx, apartmentID, entity.RuleCategory(category))
	if err != nil {
		return nil, err
	}

	resp := make([]dto.RuleResponse, len(rules))
	for i, r := range rules {
		resp[i] = *s.mapToResponse(&r)
	}
	return resp, nil
}

func (s *RuleService) Update(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID, req dto.UpdateRuleRequest) (*dto.RuleResponse, error) {
	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rule == nil || rule.ApartmentID != apartmentID {
		return nil, errors.New("rule not found")
	}

	if req.Title != "" {
		rule.Title = req.Title
	}
	if req.Description != "" {
		rule.Description = req.Description
	}
	if req.Category != "" {
		rule.Category = entity.RuleCategory(req.Category)
	}

	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return nil, err
	}

	return s.mapToResponse(rule), nil
}

func (s *RuleService) Delete(ctx context.Context, id uuid.UUID, apartmentID uuid.UUID) error {
	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rule == nil || rule.ApartmentID != apartmentID {
		return errors.New("rule not found")
	}

	return s.ruleRepo.Delete(ctx, id)
}

func (s *RuleService) mapToResponse(rule *entity.Rule) *dto.RuleResponse {
	return &dto.RuleResponse{
		ID:          rule.ID.String(),
		ApartmentID: rule.ApartmentID.String(),
		Title:       rule.Title,
		Description: rule.Description,
		Category:    string(rule.Category),
		CreatedAt:   rule.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
