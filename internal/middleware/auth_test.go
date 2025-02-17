package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthorizeRequestSucceeds(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	userID := 1
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		t.Fatalf("Could not create JWT")
	}

	c.Request.Header.Set("Cookie", "token="+tokenString)
	AuthorizeRequest(c)

	if c.IsAborted() {
		t.Fatalf("Request should not have been aborted")
	}
}

func TestAuthorizeRequestAborts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	userID := 1
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(-(time.Hour * 24)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		t.Fatalf("Could not create JWT")
	}

	c.Request.Header.Set("Cookie", "token="+tokenString)
	AuthorizeRequest(c)

	if !c.IsAborted() {
		t.Fatalf("Request should have been aborted")
	}
}
