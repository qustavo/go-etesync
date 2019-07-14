package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/andreburgaud/crypt2go/padding"
	"golang.org/x/crypto/scrypt"
)

func hmac256(salt, key []byte) []byte {
	h := hmac.New(sha256.New, salt)
	h.Write(key)
	return h.Sum(nil)
}

const blockSize = aes.BlockSize

// Cipher performs cipher operations using AES
type Cipher struct {
	cipherKey []byte
	hmacKey   []byte
}

// New returns a ne crypto object
func New(salt, key []byte) *Cipher {
	h := hmac256(salt, key)
	m := &Cipher{
		cipherKey: hmac256([]byte("aes"), h),
		hmacKey:   hmac256([]byte("hmac"), h),
	}

	return m
}

// Encrypt encrypts data
func (c *Cipher) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.cipherKey)
	if err != nil {
		return nil, err
	}

	padded, err := padding.NewPkcs7Padding(blockSize).Pad(data)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, blockSize+len(padded))
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], padded)

	return ciphertext, nil
}

// Decrypt decrypts previously encrypted data
func (c *Cipher) Decrypt(data []byte) ([]byte, error) {
	iv, ciphertext := data[:blockSize], data[blockSize:]

	block, err := aes.NewCipher(c.cipherKey)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	return padding.NewPkcs7Padding(blockSize).Unpad(plaintext)
}

// DeriveKey derives a password using scrypt
func DeriveKey(password, salt []byte) ([]byte, error) {
	return scrypt.Key(password, salt, 16384, 8, 1, 190)
}

// MustDeriveKey calls DeriveKey panicking on error
func MustDeriveKey(password, salt []byte) []byte {
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 190)
	if err != nil {
		panic(err)
	}
	return key
}
