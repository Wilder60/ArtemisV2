package domain

// KeyEntry is the struct that the request body will be serailized into
// to be passed to the business logic of the code
type KeyEntry struct {
	ID           int64  `json:"id"`
	URL          string `json:"url"`
	Name         []rune `json:"name"`
	Folder       []rune `json:"folder"`
	Username     []rune `json:"username"`
	SitePassword []rune `json:"sitepassword"`
	Notes        []rune `json:"notes"`
	Favorite     bool   `json:"favorite"`
}
