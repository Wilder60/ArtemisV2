package sql

import "github.com/Wilder60/KeyRing/internal/domain"

type SQL struct {
}

func New() SQL {
	return SQL{}
}

func (sql *SQL) GetKeyRing() []domain.KeyEntry {
	return nil
}

func (sql *SQL) AddKeyRing(domain.KeyEntry) int64 {
	return 0
}

func (sql *SQL) UpdateKeyRing(domain.KeyEntry) int64 {
	return 0
}

func (sql *SQL) DeleteKeyRing(int64) int64 {
	return 0
}
