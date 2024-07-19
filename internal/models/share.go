package models

type Share struct {
	Type     string        `json:"type"` // `file`, `directory`, `text`, `url`
	Name     string        `json:"name"`
	Password string        // hashed password
	Expiry   int           `json:"expiry"`
	Size     int64         `json:"size"`
	Files    ShareFiles    `json:"files,omitempty"`
	Creator  *ShareCreator `json:"creator,omitempty"`
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
