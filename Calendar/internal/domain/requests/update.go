package requests

import "context"

type Update struct {
	Ctx         context.Context
	UserID      string
	EventID     string `json:"eventid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	SDate       string `json:"sdate"`
	EDate       string `json:"edate"`
}

func NewUpdate(userID string) *Update {
	return &Update{
		Ctx:    context.Background(),
		UserID: userID,
	}
}
