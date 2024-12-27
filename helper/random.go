package helper

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomNumberString(length int) (string, error) {
	result := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}

		result += n.String()
	}

	return result, nil
}
