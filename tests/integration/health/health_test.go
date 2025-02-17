package health

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/FaizanAC/Go-Banking/internal/database"
	"github.com/FaizanAC/Go-Banking/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	s := server.NewServer(
		database.NewDatabase(), os.Getenv("PORT"),
	)
	r := s.SetupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}
