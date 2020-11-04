package requests

import "context"

type Delete struct {
	Ctx    context.Context
	UserID string
	ID     string `json:"ids"`
}

func NewDelete(userID string) *Delete {
	return &Delete{
		Ctx:    context.Background(),
		UserID: userID,
	}
}
