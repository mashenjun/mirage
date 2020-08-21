package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestAuth(t *testing.T) {
	openID := "ccc"
	tk, err := GenToken(openID, time.Now().Add(1*time.Hour).Unix())
	if err != nil {
		t.Fatalf("could not gen jwt token: %v", err)
	}
	t.Logf("%v", tk)

	claims, err := VerifyToken(tk)
	if err != nil {
		t.Fatalf("could not validat jwt token: %v", err)
	}
	t.Logf("%v\n", claims)
	if claims.Uid != openID {
		t.Fatal("could not validate jwt token get wrong uid")
	}
}

func TestProtection(t *testing.T) {
	openID := "ccc"
	tk, err := GenToken(openID, time.Now().Add(1*time.Hour).Unix())
	if err != nil {
		t.Fatalf("could not gen jwt token: %v", err)
	}
	t.Logf("%v", tk)
	tk = fmt.Sprintf("%sinvalid", tk)
	claims, err := VerifyToken(tk)
	if err != nil {
		t.Fatalf("could not validate jwt token: %v", err)
	}
	t.Logf("%v\n", claims)
	if claims.Uid != openID {
		t.Fatal("could not validate jwt token get wrong uid")
	}
}

func TestVerify(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJvWUNjZTVFVmQwaU9DVzI4bDNIcEJWOFB6NmVjIiwiZXhwIjoxNTk3MTA2OTgyLCJqdGkiOiJvWUNjZTVFVmQwaU9DVzI4bDNIcEJWOFB6NmVjIn0.q2tyxLrAQM2ZpPAlWi0NXb0PkuhHd8ENghuHVZjGEy0"
	claims, err := VerifyToken(token)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				t.Fatalf("expire: %+v", err)
			}
		} else {
			t.Fatalf("could not validat jwt token: %v", err)
		}
	}
	t.Logf("%v\n", claims)
}
