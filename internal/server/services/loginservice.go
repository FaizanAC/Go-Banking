package services

import (
	"fmt"

	"github.com/FaizanAC/Go-Banking/internal/models"
	"github.com/FaizanAC/Go-Banking/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginService struct {
	db *gorm.DB
}

func NewLoginService(db *gorm.DB) *LoginService {
	return &LoginService{db: db}
}

func (s *LoginService) Login(login models.Login) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", login.Email).First(&user).Error; err != nil {
		return "", fmt.Errorf("user with email %s not found", login.Email)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return "", fmt.Errorf("invalid Password")
	}

	jwtToken, err := util.GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to Generate JWT")
	}

	return jwtToken, nil
}
