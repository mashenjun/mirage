package util

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const (
	SignSecretKey = `\x88D\xf09\x6\xa0A\x7\xc5V\xbe\x8b\xef\xd7\xd8\xd3\xe6\x98*4`
)

type TokenClaims struct {
	Uid string `json:"uid"` // set as openID
	jwt.StandardClaims
}

func GenToken(uid string, expiration int64) (string, error) {
	// Create the Claims
	claims := TokenClaims{
		uid,
		jwt.StandardClaims{
			Id:        uid,
			ExpiresAt: expiration,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	st, err := token.SignedString([]byte(SignSecretKey))
	if err != nil {
		return "", err
	}

	return st, nil
}

func VerifyToken(token string) (TokenClaims, error) {
	var tokenClaims TokenClaims

	_, err := jwt.ParseWithClaims(token, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SignSecretKey), nil
	})
	if err != nil {
		return tokenClaims, err
	}
	return tokenClaims, nil
}
