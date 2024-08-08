package models

import "encoding/json"

type FileType int

const (
	FileTypeUnknown FileType = 0
	FileTypeText    FileType = iota
	FileTypeDocument
	FileTypeImage
	FileTypeVideo
	FileTypeAudio
)

func (t FileType) String() string {
	return [...]string{
		"unknown",
		"text",
		"document",
		"image",
		"video",
		"audio",
	}[t]
}

func (t FileType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
