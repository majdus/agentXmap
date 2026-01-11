package service

import (
	"agentXmap/internal/domain"
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type ResourceService interface {
	CreateResource(ctx context.Context, typeID, name string, config json.RawMessage) (*domain.Resource, error)
	GetResource(ctx context.Context, id uuid.UUID) (*domain.Resource, error)
	ListAgentsWithAccess(ctx context.Context, resourceID uuid.UUID) ([]domain.Agent, error)
}

type DefaultResourceService struct {
	resRepo domain.ResourceRepository
}

func NewResourceService(resRepo domain.ResourceRepository) *DefaultResourceService {
	return &DefaultResourceService{resRepo: resRepo}
}

func (s *DefaultResourceService) CreateResource(ctx context.Context, typeID, name string, config json.RawMessage) (*domain.Resource, error) {
	if name == "" {
		return nil, errors.New("resource name is required")
	}
	if typeID == "" {
		return nil, errors.New("resource type is required")
	}

	res := &domain.Resource{
		TypeID:            typeID,
		Name:              name,
		ConnectionDetails: config,
	}

	if err := s.resRepo.Create(ctx, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *DefaultResourceService) GetResource(ctx context.Context, id uuid.UUID) (*domain.Resource, error) {
	res, err := s.resRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("resource not found")
	}
	return res, nil
}
