package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TODO: use an interface
type Manager struct {
	key        []byte
	expireTime time.Duration
}

type Claim struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func (j *Manager) Generate(id string) (string, error) {
	d := time.Now().Add(j.expireTime)
	claim := Claim{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(d),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *Manager) Validate(t string) (*Claim, error) {
	claim := new(Claim)
	token, err := jwt.ParseWithClaims(
		t,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.key), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expired token")
	}
	return claim, nil
}

func New(key string, expireTime time.Duration) *Manager {
	return &Manager{
		key:        []byte(key),
		expireTime: expireTime,
	}
}
