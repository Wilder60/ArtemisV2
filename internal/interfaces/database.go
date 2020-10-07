package interfaces

import "github.com/Wilder60/KeyRing/internal/domain"

type Database interface {
	GetKeyRing() []domain.KeyEntry
	AddKeyRing(domain.KeyEntry) int64
	UpdateKeyRing(domain.KeyEntry) int64
	DeleteKeyRing(int64) int64
}
