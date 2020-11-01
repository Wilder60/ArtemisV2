package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/Wilder60/ShadowKeep/internal/domain"
)

type StorageMock struct {
	mock.Mock
}

func NewStorageMock() *StorageMock {
	return &StorageMock{}
}

func (s *StorageMock) GetEventsPaginated(ctx context.Context, userID, sdate string, limit, offset int) ([]domain.Event, error) {
	return []domain.Event{}, nil
}

func (s *StorageMock) GetEventsInRange(ctx context.Context, userID, sdate, edate string) ([]domain.Event, error) {
	return []domain.Event{}, nil
}
func (s *StorageMock) CreateEvents(ctx context.Context, event domain.Event) error {
	return nil
}
func (s *StorageMock) UpdateEvent(ctx context.Context, event domain.Event) error {
	return nil
}
func (s *StorageMock) DeleteEvents(ctx context.Context, userID, eventID string) error {
	return nil
}
func (s *StorageMock) DeleteEventsForUser(ctx context.Context, userID string) error {
	return nil
}
