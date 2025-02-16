package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secretKey []byte, userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)

	fmt.Println(tokenString, err)
	return tokenString, err
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected Signing Method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return token, err
}
