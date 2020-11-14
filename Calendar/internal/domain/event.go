package domain

import (
	pb "github.com/Wilder60/ArtemisV2/Calendar/pkg/grpc/service"
)

// Event is the struct that contains the superset of both the Event
// and the Reminder types
type Event struct {
	ID          string `json:"ID,omitempty"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Color       string `json:"Color"`
	SDate       string `json:"SDate"`
	EDate       string `json:"EDate"`
}

func (e *Event) ToProto() *pb.Event {
	protoEvent := &pb.Event{
		Id:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Color:       e.Color,
		Sdate:       e.SDate,
		Edate:       e.EDate,
	}
	return protoEvent
}

// time.Now().Format(time.RFC3339)
