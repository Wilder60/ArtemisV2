package domain

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

// time.Now().Format(time.RFC3339)
