package requests

import "context"

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
