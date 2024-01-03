package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateUid() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%X", b[0:]), nil
}
