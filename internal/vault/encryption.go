package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func encryptString(payload, secret string) (string, error) {
	key := []byte(secret)
	if len(key) < 32 {
		panic("32bit secret key required")
	}
	// generate AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// generate a random IV
	ciphertext := make([]byte, aes.BlockSize+len(payload))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	// encryption
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(payload))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decryptSrting(encrypted, secret string) (string, error) {
	key := []byte(secret)
	if len(key) < 32 {
		panic("32bit secret key required")
	}

	ciphertext, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Extract the IV (Initialization Vector)
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
