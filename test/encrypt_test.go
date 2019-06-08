package test

import (
	"testing"

	"github.com/sylank/lavender-commons-go/crypto"
)

const (
	PASS           = "0123456789012345"
	EMAIL          = "test@mail.hu"
	ENCRYPTED_TEXT = "Smh8lV8x5X9riRNfQFAKTaQ-TPwpa2RbacPs1A=="
)

func Test_Encryption(t *testing.T) {
	encrypted, _ := crypto.EncryptString([]byte(PASS), EMAIL)
	t.Log(encrypted)

	decrypted, _ := crypto.DecryptString([]byte(PASS), encrypted)
	t.Log(decrypted)

	if decrypted != EMAIL {
		t.Error("Decrypted data do not match with data")
		t.Fail()
	}
}
