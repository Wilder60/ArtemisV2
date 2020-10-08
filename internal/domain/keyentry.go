package domain

// KeyEntry is the struct that the request body will be serailized into
// to be passed to the business logic of the code
type KeyEntry struct {
	ID           string `json:"id"`
	UserID       string `json:"userid"`
	URL          string `json:"url"`
	SiteName     string `json:"sitename"`
	Folder       string `json:"folder"`
	Username     string `json:"username"`
	SitePassword string `json:"sitepassword"`
	Notes        string `json:"notes"`
	Favorite     bool   `json:"favorite"`
}
