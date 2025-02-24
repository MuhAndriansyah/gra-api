package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

func GenerateRandomNumberString(length int) (string, error) {
	result := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))

		if err != nil {
			return "", fmt.Errorf("failed generate random number string %v", err)
		}

		result += n.String()
	}

	return result, nil
}

func GenerateRandomString() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		return "", fmt.Errorf("failed generate random string %v", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
