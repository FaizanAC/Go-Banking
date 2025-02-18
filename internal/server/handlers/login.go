package handlers

import (
	"net/http"

	"github.com/FaizanAC/Go-Banking/internal/models"
	"github.com/FaizanAC/Go-Banking/internal/server/services"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginService *services.LoginService
}

func NewLoginHandler(loginService *services.LoginService) *LoginHandler {
	return &LoginHandler{loginService: loginService}
}

func (h *LoginHandler) HandleLogin(c *gin.Context) {
	var login models.Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := h.loginService.Login(login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", jwtToken, 3600, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login Successful"})
}
