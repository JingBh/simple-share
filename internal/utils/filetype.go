package utils

import (
	"github.com/h2non/filetype"
	"github.com/jingbh/simple-share/internal/models"
	"unicode/utf8"
)

var textPlainType = filetype.NewType("txt", "text/plain")

func textPlainMatcher(buf []byte) bool {
	return utf8.Valid(buf)
}

func IsPdf(data []byte) bool {
	return filetype.IsExtension(data, "pdf")
}

func DeduceFileType(data []byte) (models.FileType, error) {
	if filetype.IsType(data, textPlainType) {
		return models.FileTypeText, nil
	}

	if filetype.IsDocument(data) || IsPdf(data) {
		return models.FileTypeDocument, nil
	}

	if filetype.IsImage(data) {
		return models.FileTypeImage, nil
	}

	if filetype.IsAudio(data) {
		return models.FileTypeAudio, nil
	}

	if filetype.IsVideo(data) {
		return models.FileTypeVideo, nil
	}

	return models.FileTypeUnknown, nil
}

func init() {
	filetype.AddMatcher(textPlainType, textPlainMatcher)
}
