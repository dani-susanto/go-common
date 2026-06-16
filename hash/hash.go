package hash

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

type Hash interface {
	Hash(text string) (string, error)
	Compare(text, hash string) bool
	GenerateToken(length int) (string, error)
	GenerateOTP(length int) (string, error)
}

type hash struct{}

func (hash) Hash(text string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (hash) Compare(text, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text)) == nil
}

func (hash) GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (h *hash) GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}
	return string(otp), nil
}

func New() Hash {
	return &hash{}
}
