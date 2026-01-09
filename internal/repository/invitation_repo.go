package repository

import (
	"context"

	"agentXmap/internal/domain"

	"gorm.io/gorm"
)

type invitationRepository struct {
	db *gorm.DB
}

// NewInvitationRepository creates a new postgres repository for Invitations.
func NewInvitationRepository(db *gorm.DB) domain.InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) Create(ctx context.Context, invitation *domain.Invitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *invitationRepository) GetByToken(ctx context.Context, token string) (*domain.Invitation, error) {
	var invitation domain.Invitation
	if err := r.db.WithContext(ctx).
		Preload("Organization").
		Preload("Invitor").
		First(&invitation, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *invitationRepository) Update(ctx context.Context, invitation *domain.Invitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}
