package interfaces

import (
	"github.com/Wilder60/KeyRing/internal/domain"
)

type Database interface {
	GetKeyRing(string, int64, int64) ([]domain.KeyEntry, error)
	AddKeyRing(domain.KeyEntry) (int64, error)
	UpdateKeyRing(domain.KeyEntry) (int64, error)
	DeleteKeyRing(string) (int64, error)
}
