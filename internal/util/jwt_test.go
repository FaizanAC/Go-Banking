package util

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseJWT(t *testing.T) {
	var userID uint = 1
	tokenString, err := GenerateJWT(userID)
	if err != nil {
		t.Fatalf("GenerateJWT threw an Error %v", err)
	}

	jwtToken, err := ParseJWT(tokenString)
	if err != nil {
		t.Fatalf("ParseJWT failed to Parse %v", err)
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		if claims["sub"] != float64(userID) {
			t.Fatalf("JWT UserID does not match")
		}
	}
}
