package service

import (
	"agentXmap/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuditRepository
type MockAuditRepository struct {
	mock.Mock
}

func (m *MockAuditRepository) CreateLog(ctx context.Context, log *domain.SystemAuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockAuditRepository) CreateExecution(ctx context.Context, exec *domain.AgentExecution) error {
	args := m.Called(ctx, exec)
	return args.Error(0)
}

func TestAuditService_LogAction(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	entityID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAuditRepository)
		service := NewAuditService(mockRepo)

		mockRepo.On("CreateLog", ctx, mock.MatchedBy(func(l *domain.SystemAuditLog) bool {
			return l.Action == domain.AuditActionCreate && l.EntityType == "agent"
		})).Return(nil)

		err := service.LogAction(ctx, &userID, "agent", entityID, domain.AuditActionCreate, json.RawMessage(`{}`), "127.0.0.1")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockAuditRepository)
		service := NewAuditService(mockRepo)

		mockRepo.On("CreateLog", ctx, mock.Anything).Return(errors.New("db error"))

		err := service.LogAction(ctx, &userID, "agent", entityID, domain.AuditActionCreate, nil, "")
		assert.Error(t, err)
		if err != nil {
			assert.Equal(t, "db error", err.Error())
		}
	})
}

func TestAuditService_RecordExecution(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAuditRepository)
	service := NewAuditService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		exec := &domain.AgentExecution{CreatedAt: time.Now()}
		mockRepo.On("CreateExecution", ctx, exec).Return(nil)

		err := service.RecordExecution(ctx, exec)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
