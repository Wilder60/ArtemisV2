package requests

import "context"

type GetRange struct {
	Ctx    context.Context
	UserID string
	Sdate  string `form:"sdate"`
	Edate  string `form:"edate"`
}

func NewGetRange(userID string) *GetRange {
	return &GetRange{
		Ctx:    context.Background(),
		UserID: userID,
	}
}
