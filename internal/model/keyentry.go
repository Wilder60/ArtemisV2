package model

// Event is the struct that the request body will be serailized into
// to be passed to the business logic of the code
type KeyEntry struct {
	URL          string
	Name         string
	Folder       string
	Username     string
	SitePassword string
	Notes        string
	Favorite     bool
}
