package requests

import (
	"context"

	pb "github.com/Wilder60/ArtemisV2/Calendar/pkg/grpc/service"
)

type GetRange struct {
	Ctx    context.Context
	UserID string
	Sdate  string `form:"sdate"`
	Edate  string `form:"edate"`
}

func NewProtoGetRange(userID string, req *pb.GetEventsInRangeRequest) *GetRange {
	return &GetRange{
		Ctx:    context.Background(),
		UserID: userID,
		Sdate:  req.Sdate,
		Edate:  req.Edate,
	}
}

func NewGetRange(userID string) *GetRange {
	return &GetRange{
		Ctx:    context.Background(),
		UserID: userID,
	}
}
