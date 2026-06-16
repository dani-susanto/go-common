package jwt

import (
	"errors"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT interface {
	Generate(data *ClaimsInput) (*ClaimsOutput, error)
	Validate(tokenString string) (*ClaimsInput, error)
}

type jwt struct {
	secret    []byte
	expiredIn time.Duration
}

type contextKey string

const ClaimsKey contextKey = "claims"

type ClaimsInput struct {
	UserID string `json:"user_id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
	gojwt.RegisteredClaims
}

type ClaimsOutput struct {
	ID        string
	Token     string
	IssuedAt  *gojwt.NumericDate
	ExpiresAt *gojwt.NumericDate
}

func New(secret string, expiredIn time.Duration) JWT {
	return &jwt{
		secret:    []byte(secret),
		expiredIn: expiredIn,
	}
}

func (j *jwt) Generate(data *ClaimsInput) (*ClaimsOutput, error) {
	now := time.Now()

	data.RegisteredClaims = gojwt.RegisteredClaims{
		Subject:   data.UserID,
		Issuer:    "account-service",
		Audience:  []string{"public"},
		IssuedAt:  gojwt.NewNumericDate(now),
		ExpiresAt: gojwt.NewNumericDate(now.Add(j.expiredIn)),
		ID:        data.RegisteredClaims.ID,
	}

	if len(data.RegisteredClaims.ID) == 0 {
		data.RegisteredClaims.ID = fmt.Sprintf("%s-%s", data.UserID, uuid.New().String())
	}

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, data).SignedString(j.secret)
	if err != nil {
		return nil, err
	}

	return &ClaimsOutput{
		ID:        data.RegisteredClaims.ID,
		Token:     token,
		IssuedAt:  data.RegisteredClaims.IssuedAt,
		ExpiresAt: data.RegisteredClaims.ExpiresAt,
	}, nil
}

func (j *jwt) Validate(tokenString string) (*ClaimsInput, error) {
	token, err := gojwt.ParseWithClaims(tokenString, &ClaimsInput{}, func(token *gojwt.Token) (interface{}, error) {
		if token.Method != gojwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClaimsInput)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
