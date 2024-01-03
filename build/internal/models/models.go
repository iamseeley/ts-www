package models

// Content represents the basic fields required for rendering content on the site
type Page struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        []byte `json:"body"`
	Theme       string `json:"theme"` // Assuming theme is consistent across content
	Collection  string `json:"collection"`
}
