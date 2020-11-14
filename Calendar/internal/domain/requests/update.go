package requests

import (
	"context"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain"
	pb "github.com/Wilder60/ArtemisV2/Calendar/pkg/grpc/service"
)

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

func NewProtoUpdate(userID string, req *pb.UpdateEventRequest) *Update {
	return &Update{
		Ctx:         context.Background(),
		UserID:      userID,
		EventID:     req.Eventid,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		SDate:       req.Sdate,
		EDate:       req.Edate,
	}
}

func NewUpdate(userID string) *Update {
	return &Update{
		Ctx:    context.Background(),
		UserID: userID,
	}
}

func (up *Update) ToEvent() *domain.Event {
	return &domain.Event{
		ID:          up.EventID,
		Name:        up.Name,
		Description: up.Description,
		Color:       up.Color,
		SDate:       up.SDate,
		EDate:       up.EDate,
	}
}
