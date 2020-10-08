package interfaces

import (
	"github.com/Wilder60/KeyRing/internal/domain"
	"github.com/google/uuid"
)

type Database interface {
	GetKeyRing(uuid.UUID, int64, int64) ([]domain.KeyEntry, error)
	AddKeyRing(domain.KeyEntry) (int64, error)
	UpdateKeyRing(domain.KeyEntry) (int64, error)
	DeleteKeyRing(int64) (int64, error)
}
