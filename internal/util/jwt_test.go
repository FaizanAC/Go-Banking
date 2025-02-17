package util

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseJWT(t *testing.T) {
	var userID uint = 1
	tokenString, err := GenerateJWT(userID)
	assert.Nil(t, err)

	jwtToken, err := ParseJWT(tokenString)
	assert.Nil(t, err)

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, claims["sub"], float64(userID))
	}
}
