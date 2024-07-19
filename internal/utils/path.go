package utils

import "strings"

func ExtractFilename(path string) string {
	path = strings.TrimSpace(path)
	path = strings.TrimSuffix(path, "/")
	parts := strings.SplitAfter(path, "/")
	return parts[len(parts)-1]
}
