package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateSignedToken(tkn *jwt.Token, p8key []byte) (string, error) {
	key, err := readPrivateKey(p8key)
	if err != nil {
		return "", err
	}

	signed, err := tkn.SignedString(key)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func GenerateToken(issuer, kid string) *jwt.Token {
	tkn := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		jwt.StandardClaims{
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
		},
	)
	tkn.Header["kid"] = kid
	return tkn
}

func readPrivateKey(p8key []byte) (*ecdsa.PrivateKey, error) {
	// Here you need to decode the Apple private key, which is in pem format
	block, _ := pem.Decode(p8key)
	// Check if it's a private key
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// Now you need an instance of *ecdsa.PrivateKey
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("not ecdsa private key")
	}

	return ecdsaPrivateKey, nil
}
