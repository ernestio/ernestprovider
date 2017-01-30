package ernestAes

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "lkasdisiucnlisucnliausui"
	plaintext := "foo"
	c := New()

	cipher, err := c.Encrypt(plaintext, key)
	if err != nil {
		t.Errorf("An unexpected error occured")
	}
	processed, err := c.Decrypt(cipher, key)
	if err != nil {
		t.Errorf("An unexpected error occured")
	}
	a := string(processed)
	b := string(plaintext)

	if a != b {
		t.Errorf("Processed string is not equal as initial plaintext")
		return
	}

}
