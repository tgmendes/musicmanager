package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type Apple struct {
	iss        string
	kid        string
	privateKey *ecdsa.PrivateKey
	token      *jwt.Token
}

func NewApple(issuer, kid string, p8key []byte) (*Apple, error) {
	a := Apple{
		iss: issuer,
		kid: kid,
	}

	pKey, err := readPrivateKey(p8key)
	if err != nil {
		return nil, err
	}
	a.privateKey = pKey

	a.token = a.GenerateToken()
	return &a, nil
}

func (a *Apple) SignedToken() (string, error) {
	if a.token == nil {
		a.token = a.GenerateToken()
	}

	if err := a.token.Claims.Valid(); err != nil {
		a.token = a.GenerateToken()
	}
	signed, err := a.token.SignedString(a.privateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (a *Apple) GenerateToken() *jwt.Token {
	tkn := jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": jwt.SigningMethodES256.Alg(),
			"kid": a.kid,
		},
		Claims: jwt.StandardClaims{
			Issuer:    a.iss,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
		},
		Method: jwt.SigningMethodES256,
	}

	return &tkn
}

func readPrivateKey(key []byte) (*ecdsa.PrivateKey, error) {
	// Here you need to decode the Apple private key, which is in pem format
	block, _ := pem.Decode(key)
	// Check if it's a private key
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// Now you need an instance of *ecdsa.PrivateKey
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("not ecdsa private key")
	}

	return ecdsaPrivateKey, nil
}
