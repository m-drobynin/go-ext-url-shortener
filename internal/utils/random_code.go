package utils

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const alphabetLength int64 = int64(len(alphabet))

func RandomCode(length int) (string, error) {
	code := make([]byte, length)

	for i := range code {
		index, err := rand.Int(rand.Reader, big.NewInt(alphabetLength))
		if err != nil {
			return "", err
		}

		code[i] = alphabet[index.Int64()]
	}

	return string(code), nil
}
