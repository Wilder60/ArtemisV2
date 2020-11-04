package adapter

import (
	"context"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain/requests"
)

// Storage is the interface that is used by the adapter
type Storage interface {
	GetEventsPaginated(context.Context, requests.GetPagination) ([]domain.Event, error)
	GetEventsInRange(context.Context, requests.GetRange) ([]domain.Event, error)
	CreateEvents(context.Context, requests.Add) error
	UpdateEvent(context.Context, requests.Update) error
	DeleteEvents(context.Context, requests.Delete) error
	DeleteEventsForUser(context.Context, requests.Delete) error
}
