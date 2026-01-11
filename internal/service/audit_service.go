package service

import (
	"agentXmap/internal/domain"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type AuditService interface {
	LogAction(ctx context.Context, actorUserID *uuid.UUID, entityType string, entityID uuid.UUID, action domain.AuditAction, changes json.RawMessage, ipAddress string) error
	RecordExecution(ctx context.Context, exec *domain.AgentExecution) error
}

type DefaultAuditService struct {
	auditRepo domain.AuditRepository
}

func NewAuditService(auditRepo domain.AuditRepository) *DefaultAuditService {
	return &DefaultAuditService{auditRepo: auditRepo}
}

func (s *DefaultAuditService) LogAction(ctx context.Context, actorUserID *uuid.UUID, entityType string, entityID uuid.UUID, action domain.AuditAction, changes json.RawMessage, ipAddress string) error {
	log := &domain.SystemAuditLog{
		ActorUserID: actorUserID,
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      action,
		Changes:     changes,
		IPAddress:   ipAddress,
	}
	return s.auditRepo.CreateLog(ctx, log)
}

func (s *DefaultAuditService) RecordExecution(ctx context.Context, exec *domain.AgentExecution) error {
	// Here we could add logic to anonymize or calculate missing metrics if needed.
	// For now, it's a direct record.
	return s.auditRepo.CreateExecution(ctx, exec)
}
