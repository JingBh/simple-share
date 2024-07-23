package oss

import (
	"fmt"
	"math/rand"
	"regexp"
)

// excluded 0, 1, l, o to avoid confusion
const charset = "abcdefghijkmnpqrstuvwxyz23456789"

// allowed name pattern
const validPattern = `^[a-zA-Z0-9]{2,64}$`

func generateShareName(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateShareName(length int) (string, error) {
	for i := 0; i < 5; i++ {
		name := generateShareName(length)
		if CheckShareName(name) {
			return name, nil
		}
	}
	return "", fmt.Errorf("unable to generate a unique share name")
}

func CheckShareName(name string) bool {
	// check pattern
	re, err := regexp.Compile(validPattern)
	if err != nil || !re.MatchString(name) {
		return false
	}

	// check existence
	client := Client()
	_, err = client.GetObjectMeta("shares/" + name)
	return err != nil // object does not exist, then name is available
}
