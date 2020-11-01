package adapter

import (
	"context"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain"
)

type Storage interface {
	GetEventsPaginated(context.Context, string, string, int, int) ([]domain.Event, error)
	GetEventsInRange(context.Context, string, string, string) ([]domain.Event, error)
	CreateEvents(context.Context, domain.Event) error
	UpdateEvent(context.Context, domain.Event) error
	DeleteEvents(context.Context, string, string) error
	DeleteEventsForUser(context.Context, string) error
}
