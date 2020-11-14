package requests

import (
	"context"

	pb "github.com/Wilder60/ArtemisV2/Calendar/pkg/grpc/service"
)

type Delete struct {
	Ctx    context.Context
	UserID string
	ID     string `json:"ids"`
}

func NewProtoDelete(userID string, req *pb.DeleteRequest) *Delete {
	return &Delete{
		Ctx:    context.Background(),
		UserID: userID,
		ID:     req.Id,
	}
}

func NewDelete(userID string) *Delete {
	return &Delete{
		Ctx:    context.Background(),
		UserID: userID,
	}
}
