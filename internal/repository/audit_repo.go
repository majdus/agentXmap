package repository

import (
	"context"

	"agentXmap/internal/domain"

	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) domain.AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) CreateLog(ctx context.Context, log *domain.SystemAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditRepository) CreateExecution(ctx context.Context, exec *domain.AgentExecution) error {
	// For partitioned tables, simple Insert usually works if specific partitions exist.
	return r.db.WithContext(ctx).Create(exec).Error
}
