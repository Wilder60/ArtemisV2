package requests

import (
	"context"

	pb "github.com/Wilder60/ArtemisV2/Calendar/pkg/grpc/service"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain"
	"github.com/gofrs/uuid"
)

// Add struct that will be passed to the the storage object
// This contains information that will be pulled from the
type Add struct {
	Ctx         context.Context
	UserID      string
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	SDate       string `json:"sdate"`
	EDate       string `json:"edate"`
}

// NewAdd is the constructor for the Add struct, this will initalize a new context for the
// firestore request.  This function takes the userID as a parameters that will be passed
// from the gin context
func NewAdd(userID string) *Add {
	return &Add{
		Ctx:    context.Background(),
		UserID: userID,
	}
}

// BindProto will take a proto AddEventRequest and bind the data to the add request object
// then the object can be passed to the db for processing
func NewProtoAdd(userID string, req *pb.AddEventRequest) *Add {
	return &Add{
		Ctx:         context.Background(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		SDate:       req.Sdate,
		EDate:       req.Edate,
	}
}

// ToEvent will generate a new uuid for the event, and store it in the Event object
func (add *Add) ToEvent() (*domain.Event, error) {
	eventID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &domain.Event{
		ID:          eventID.String(),
		Name:        add.Name,
		Description: add.Description,
		Color:       add.Color,
		SDate:       add.SDate,
		EDate:       add.EDate,
	}, nil
}
