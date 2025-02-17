package server

import (
	"net/http"
	"time"

	"github.com/FaizanAC/Go-Banking/internal/models"
	"github.com/FaizanAC/Go-Banking/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
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

	jwtToken, err := util.GenerateJWT(user.ID)
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
	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	newAccount := models.BankAccount{
		AccountNumber: util.GenerateAccountNumber(),
		UserID:        userID.(uint),
		Balance:       0,
	}

	if res := s.db.Create(&newAccount); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Create Account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Account Number": newAccount.AccountNumber})
}

func (s *Server) handleGetAccount(c *gin.Context) {
	var account models.BankAccount
	id := c.Param("id")

	if err := s.db.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not Found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s *Server) handleDeposit(c *gin.Context) {
	var deposit models.Transaction

	if err := c.ShouldBindJSON(&deposit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var account models.BankAccount
	if err := s.db.Where("account_number = ?", deposit.AccountNumber).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not Found"})
		return
	}

	if account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this Account"})
		return
	}

	account.Balance += deposit.Amount

	deposit.Type = "DEPOSIT"
	deposit.TransactionID = uuid.New().String()

	eg := errgroup.Group{}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		eg.Go(func() error {
			return tx.Save(&account).Error
		})

		eg.Go(func() error {
			return tx.Save(&deposit).Error
		})

		return eg.Wait()
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Deposit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": account.Balance})
}

func (s *Server) handleWithdraw(c *gin.Context) {
	var withdraw models.Transaction

	if err := c.ShouldBindJSON(&withdraw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var account models.BankAccount
	if err := s.db.Where("account_number = ?", withdraw.AccountNumber).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not Found"})
		return
	}

	if account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this Account"})
		return
	}

	if account.Balance < withdraw.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient Balance"})
		return
	}

	account.Balance -= withdraw.Amount

	withdraw.Type = "WITHDRAW"
	withdraw.TransactionID = uuid.New().String()

	eg := errgroup.Group{}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		eg.Go(func() error {
			return tx.Save(&account).Error
		})

		eg.Go(func() error {
			return tx.Save(&withdraw).Error
		})

		return eg.Wait()
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Withdrawl"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": account.Balance})
}

func (s *Server) handleSendTransfer(c *gin.Context) {
	var transfer struct {
		Amount        float64 `json:"amount" binding:"required"`
		AccountNumber string  `json:"accountNumber" binding:"required"`
		ReceiverID    uint    `json:"receiverID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var senderAccount models.BankAccount
	if err := s.db.Where("account_number = ?", transfer.AccountNumber).First(&senderAccount).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender Account not Found"})
		return
	}

	if senderAccount.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this Account"})
		return
	}

	transactionDetails := models.Transaction{
		Amount:        transfer.Amount,
		AccountNumber: transfer.AccountNumber,
		TransactionID: uuid.New().String(),
		Type:          "TRANSFER",
	}

	transferRow := models.Transfer{
		SenderID:      userID.(uint),
		ReceiverID:    transfer.ReceiverID,
		Amount:        transfer.Amount,
		Status:        "PENDING",
		ExpiresOn:     time.Now().Add(time.Second * 3600 * 30),
		TransactionID: transactionDetails.TransactionID,
	}

	senderAccount.Balance -= transfer.Amount

	eg := errgroup.Group{}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		eg.Go(func() error {
			return tx.Save(&senderAccount).Error
		})

		eg.Go(func() error {
			return tx.Save(&transferRow).Error
		})

		eg.Go(func() error {
			return tx.Save(&transactionDetails).Error
		})

		return eg.Wait()
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Send Transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": senderAccount.Balance})
}

func (s *Server) handleAcceptTransfer(c *gin.Context) {
	var acceptTransfer struct {
		TransactionID string `json:"transactionID" binding:"required"`
		AccountNumber string `json:"accountNumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&acceptTransfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var tranferDetails models.Transfer
	if err := s.db.Where("transaction_id = ?", acceptTransfer.TransactionID).First(&tranferDetails).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender Account not Found"})
		return
	}

	if tranferDetails.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not able to accept this Trasfer"})
		return
	}

	var userAccount models.BankAccount
	if err := s.db.Where("account_number = ?", acceptTransfer.AccountNumber).Find(&userAccount).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "The desired Account does not exist"})
		return
	}

	userAccount.Balance += tranferDetails.Amount
	tranferDetails.Status = "ACCEPTED"

	transactionDetails := models.Transaction{
		Amount:        tranferDetails.Amount,
		AccountNumber: userAccount.AccountNumber,
		TransactionID: uuid.New().String(),
		Type:          "TRANSFER",
	}

	eg := errgroup.Group{}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		eg.Go(func() error {
			return tx.Save(&userAccount).Error
		})

		eg.Go(func() error {
			return tx.Save(&tranferDetails).Error
		})

		eg.Go(func() error {
			return tx.Save(&transactionDetails).Error
		})

		return eg.Wait()
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Accept Transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": userAccount.Balance})
}

func (s *Server) handleActivityFeed(c *gin.Context) {
	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var allAccounts []models.BankAccount
	if err := s.db.Where("user_id = ?", userID).Find(&allAccounts).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"lastestActivity": nil})
		return
	}

	accountNumbers, seenNumber := []string{}, make(map[string]bool)
	for _, nextAccount := range allAccounts {
		if _, hasKey := seenNumber[nextAccount.AccountNumber]; !hasKey {
			accountNumbers = append(accountNumbers, nextAccount.AccountNumber)
			seenNumber[nextAccount.AccountNumber] = true
		}
	}

	var latestTransactions []models.Transaction
	if err := s.db.Where("account_number IN ?", accountNumbers).Order("created_at desc").Limit(10).Find(&latestTransactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get Activity Feed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"latestActivity": latestTransactions})
}
