package models

import "time"

type Share struct {
	Type      string        `json:"type"` // `file`, `directory`, `text`, `url`
	Name      string        `json:"name"`
	Password  string        `json:"password,omitempty"` // hashed password
	Expiry    int           `json:"expiry,omitempty"`
	Size      int64         `json:"size"`
	CreatedAt time.Time     `json:"createdAt,omitempty"`
	Files     ShareFiles    `json:"files,omitempty"`
	Creator   *ShareCreator `json:"creator,omitempty"`
}

type ShareFiles []struct {
	Id   string `json:"id"`
	Path string `json:"path"`
	Size int64  `json:"size,omitempty"`
}

type ShareCreator struct {
	Subject  string `json:"subject"`
	Username string `json:"username,omitempty"`
}
