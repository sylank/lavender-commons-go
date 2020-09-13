package crypto

import (
	"fmt"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	testCases := []struct {
		desc          string
		plainText     string
		encryptedText string
		key           string
	}{
		{
			desc:          "Plain text after encryption and decryption equals",
			plainText:     "test message: Hello World!!1!",
			encryptedText: "+WFwMBKLc0fxmoOY4c5wsRQtH88YpM8Tjf586dErOVFi6X/HClMk+wkUKN+P",
			key:           "12345678901234567890123456789012",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			encrypted, err := EncryptString([]byte(tC.key), tC.plainText)
			t.Log(fmt.Sprintf("encrypted text: %s", encrypted))

			if err != nil {
				t.Log(err)
				t.Fail()
			}

			if encrypted != tC.encryptedText {
				t.Fail()
			}

			decrypted, err := DecryptString([]byte(tC.key), encrypted)
			t.Log(fmt.Sprintf("decrypted text: %s", decrypted))

			if err != nil {
				t.Log(err)
				t.Fail()
			}

			if decrypted != tC.plainText {
				t.Fail()
			}
		})
	}
}
