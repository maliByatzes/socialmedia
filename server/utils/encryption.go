package utils

import (
	"crypto/aes"
	"fmt"
	"os"
)

func encryptData(data []byte) (string, error) {
	key, ok := os.LookupEnv("CRYTO_KEY")
	if !ok {
		panic("CRYPTO_KEY is not set!")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to create new cipher: %v", err)
	}

	var encryptedData []byte

	block.Encrypt(encryptedData, data)

	return string(encryptedData), nil
}

func decryptData(data []byte) (string, error) {
	key, ok := os.LookupEnv("CRYTO_KEY")
	if !ok {
		panic("CRYPTO_KEY is not set!")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to create new cipher: %v", err)
	}

	var decryptedData []byte

	block.Decrypt(decryptedData, data)

	return string(decryptedData), nil
}
