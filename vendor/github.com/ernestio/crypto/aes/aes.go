package ernestAes

import (
	"encoding/base64"
	"errors"
	"io"
	"log"

	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"

	"golang.org/x/crypto/pbkdf2"
)

const saltlen = 8
const keylen = 32
const iterations = 10002

// Crypto : ...
type Crypto struct {
}

// New : Constructor for aes.Crypto
func New() Crypto {
	return Crypto{}
}

// KeyValidation : Checks if a key is valid or not
func (aesCrypto Crypto) KeyValidation(key []byte) bool {
	l := len(key)
	if l == 16 || l == 32 || l == 64 {
		return true
	}
	return false
}

// Encrypt : Encrypt a string based on a key
func (aesCrypto Crypto) Encrypt(plaintext, password string) (string, error) {
	header := make([]byte, saltlen+aes.BlockSize)

	salt := header[:saltlen]
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		log.Println(err.Error())
		return plaintext, err
	}

	iv := header[saltlen : aes.BlockSize+saltlen]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println(err.Error())
		return plaintext, err
	}

	key := pbkdf2.Key([]byte(password), salt, iterations, keylen, md5.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
		return plaintext, err
	}

	ciphertext := make([]byte, len(header)+len(plaintext))
	copy(ciphertext, header)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize+saltlen:], []byte(plaintext))

	return base64Encode(ciphertext), nil

}

// Decrypt : Decripts a cipher text based on a key
func (aesCrypto Crypto) Decrypt(encrypted, password string) (string, error) {
	a, err := base64Decode([]byte(encrypted))
	if err != nil {
		log.Println(err.Error())
		return encrypted, err
	}
	ciphertext := a
	salt := ciphertext[:saltlen]
	iv := ciphertext[saltlen : aes.BlockSize+saltlen]
	key := pbkdf2.Key([]byte(password), salt, iterations, keylen, md5.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
		return encrypted, err
	}

	if len(ciphertext) < aes.BlockSize {
		log.Println("Invalid ciphertext")
		return "", errors.New("Invalid ciphertext")
	}

	decrypted := ciphertext[saltlen+aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrypted, decrypted)

	return string(decrypted), nil
}

func base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
