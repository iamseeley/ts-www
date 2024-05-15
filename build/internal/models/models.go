package models

type Content struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Body            []byte `json:"body"`
	Draft           bool   `json:"draft"`
	URL             string `json:"URL"`
	Featured        bool   `json:"featured,omitempty"`
	Theme           string `json:"theme"`
	Collection      string `json:"collection"`
	Date            string `json:"date,omitempty"`
	DataTitle       string `json:"data-title,omitempty"`
	DataDescription string `json:"data-description,omitempty"`
	DataImage       string `json:"data-image,omitempty"`
}
