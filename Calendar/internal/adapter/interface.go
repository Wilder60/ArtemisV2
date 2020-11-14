package adapter

import (
	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain/requests"
)

// Storage is the interface that is used by the adapter
type Storage interface {
	GetEventsPaginated(*requests.GetPagination) ([]domain.Event, error)
	GetEventsInRange(*requests.GetRange) ([]domain.Event, error)
	CreateEvents(*requests.Add) error
	UpdateEvent(*requests.Update) error
	DeleteEvents(*requests.Delete) error
	DeleteEventsForUser(*requests.Delete) error
}
