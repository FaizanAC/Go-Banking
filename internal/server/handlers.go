package server

import (
	"net/http"
	"os"

	"github.com/FaizanAC/Go-Banking/internal/models"
	"github.com/FaizanAC/Go-Banking/internal/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (s *Server) handleLogin(c *gin.Context) {
	var login struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := s.db.Where("email = ?", login.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not Found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Password"})
		return
	}

	jwtToken, err := util.GenerateJWT([]byte(os.Getenv("JWT_KEY")), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Generate JWT"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", jwtToken, 3600, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login Successful"})
}

func (s *Server) handleGetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := s.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not Found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) handleUserCreation(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Encrypt Password"})
		return
	}

	user.Password = string(encryptedPassword)

	if res := s.db.Create(&user); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Create User"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"email": user.Email})
}

func (s *Server) handleNewAccount(c *gin.Context) {
	var account models.BankAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	account.UserID = userID.(uint)

	if res := s.db.Create(&account); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Create Account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Account Number": account.AccountNumber})
}

func (s *Server) handleGetAccount(c *gin.Context) {
	var account models.BankAccount
	id := c.Param("id")

	if err := s.db.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not Found"})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	account.UserID = userID.(uint)

	c.JSON(http.StatusOK, account)
}

func (s *Server) handleDeposit(c *gin.Context) {
	var deposit struct {
		AccountNumber string  `json:"account_number" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&deposit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account models.BankAccount
	if err := s.db.Where("account_number = ?", deposit.AccountNumber).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not Found"})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	account.UserID = userID.(uint)

	account.Balance += deposit.Amount

	if res := s.db.Save(&account); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Deposit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": account.Balance})
}

func (s *Server) handleWithdraw(c *gin.Context) {
	var withdraw struct {
		AccountNumber string  `json:"account_number" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&withdraw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account models.BankAccount
	if err := s.db.Where("account_number = ?", withdraw.AccountNumber).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not Found"})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	account.UserID = userID.(uint)

	if account.Balance < withdraw.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient Balance"})
		return
	}

	account.Balance -= withdraw.Amount

	if res := s.db.Save(&account); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Withdraw"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": account.Balance})
}

func (s *Server) handleTransfer(c *gin.Context) {
	var transfer struct {
		SenderAccountNumber   string  `json:"sender_account_number" binding:"required"`
		ReceiverAccountNumber string  `json:"receiver_account_number" binding:"required"`
		Amount                float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var senderAccount models.BankAccount
	if err := s.db.Where("account_number = ?", transfer.SenderAccountNumber).First(&senderAccount).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender Account not Found"})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	senderAccount.UserID = userID.(uint)

	if senderAccount.Balance < transfer.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient Balance"})
		return
	}

	var receiverAccount models.BankAccount
	if err := s.db.Where("account_number = ?", transfer.ReceiverAccountNumber).First(&receiverAccount).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver Account not Found"})
		return
	}

	senderAccount.Balance -= transfer.Amount
	receiverAccount.Balance += transfer.Amount

	tx := s.db.Begin()
	if res := tx.Save(&senderAccount); res.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Transfer"})
		return
	}

	if res := tx.Save(&receiverAccount); res.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Transfer"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"New Balance": senderAccount.Balance})
}
