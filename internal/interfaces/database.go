package interfaces

import (
	"github.com/Wilder60/KeyRing/internal/domain"
)

// Database is contains the functions to be used by the service this allows the
// service to use any underlying database and not having to change the driving code
type Database interface {
	GetKeyRing(string, int64, int64) ([]domain.KeyEntry, error)
	AddKeyRing(domain.KeyEntry, string) (int64, error)
	UpdateKeyRing(domain.KeyEntry, string) (int64, error)
	DeleteKeyRing(string) (int64, error)
}
