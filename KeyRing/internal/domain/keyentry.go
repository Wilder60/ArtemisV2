package domain

// KeyEntry is the struct that the request body will be serailized into
// to be passed to the business logic of the code
type KeyEntry struct {
	ID           string `json:"id,omitempty"`
	URL          string `json:"url,omitempty"`
	SiteName     string `json:"sitename,omitempty"`
	Folder       string `json:"folder,omitempty"`
	Username     string `json:"username,omitempty"`
	SitePassword string `json:"sitepassword,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Favorite     bool   `json:"favorite,omitempty"`
}
