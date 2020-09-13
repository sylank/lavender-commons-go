package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const noncesize = 12

// EncryptString ...
func EncryptString(key []byte, input string) (string, error) {
	plaintext := []byte(input)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, noncesize)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString ...
func DecryptString(key []byte, securemess string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(securemess)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, noncesize)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decrypted, err := aesgcm.Open(nil, nonce, []byte(ciphertext), nil)
	return string(decrypted), err
}
