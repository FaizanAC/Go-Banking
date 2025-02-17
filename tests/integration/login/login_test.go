package login

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/FaizanAC/Go-Banking/internal/database"
	"github.com/FaizanAC/Go-Banking/internal/models"
	"github.com/FaizanAC/Go-Banking/internal/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginWithValidAccount(t *testing.T) {
	userPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 10)
	assert.Nil(t, err)

	db := database.NewDatabase()
	database.MigrateDB(db)
	res := db.Create(&models.User{
		Email:    "test@example.com",
		Password: string(userPassword),
	})
	assert.Nil(t, res.Error)

	s := server.NewServer(
		database.NewDatabase(), os.Getenv("PORT"),
	)
	r := s.SetupRouter()

	w := httptest.NewRecorder()

	jsonBody := `{"email": "test@example.com", "password": "password"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login Successful")
}

func TestLoginWithInvalidAccount(t *testing.T) {
	s := server.NewServer(
		database.NewDatabase(), os.Getenv("PORT"),
	)
	r := s.SetupRouter()

	w := httptest.NewRecorder()

	jsonBody := `{"email": "bad@example.com", "password": "password"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(jsonBody))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not Found")
}
