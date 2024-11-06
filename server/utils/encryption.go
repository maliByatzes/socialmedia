package utils

import (
	"crypto/aes"
	"log"
	"os"
)

func EncryptData(data []byte) string {
	key, ok := os.LookupEnv("CRYTO_KEY")
	if !ok {
		panic("CRYPTO_KEY is not set!")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		// NOTE: panic here for some reason
		log.Println("%w", err)
		panic("failed to create new cipher")
	}

	var encryptedData []byte

	block.Encrypt(encryptedData, data)

	return string(encryptedData)
}

func DecryptData(data []byte) string {
	key, ok := os.LookupEnv("CRYTO_KEY")
	if !ok {
		panic("CRYPTO_KEY is not set!")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		// NOTE: panic here for some reason
		log.Println("%w", err)
		panic("failed to create new cipher")
	}

	var decryptedData []byte

	block.Decrypt(decryptedData, data)

	return string(decryptedData)
}
