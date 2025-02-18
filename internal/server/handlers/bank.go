package handlers

import (
	"net/http"

	"github.com/FaizanAC/Go-Banking/internal/models"
	"github.com/FaizanAC/Go-Banking/internal/server/services"
	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	bankService *services.BankService
}

func NewBankHandler(bankService *services.BankService) *BankHandler {
	return &BankHandler{bankService}
}

func (h *BankHandler) HandleNewAccount(c *gin.Context) {
	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	newAccount, err := h.bankService.CreateAccount(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Account Number": newAccount.AccountNumber})
}

func (s *BankHandler) HandleGetAccounts(c *gin.Context) {
	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	allAccounts, err := s.bankService.GetAccountsByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allAccounts)
}

func (s *BankHandler) HandleDeposit(c *gin.Context) {
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

	account, err := s.bankService.DepositToAccount(deposit, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": account.Balance})
}

func (s *BankHandler) HandleWithdraw(c *gin.Context) {
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

	account, err := s.bankService.WithdrawFromAccount(withdraw, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": account.Balance})
}

func (s *BankHandler) HandleActivityFeed(c *gin.Context) {
	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	latestTransactions, err := s.bankService.GetActivityFeed(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"latestActivity": latestTransactions})
}

func (s *BankHandler) HandleSendTransfer(c *gin.Context) {
	var transfer models.OutgoingTransfer

	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	senderAccount, err := s.bankService.SendTransfer(transfer, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": senderAccount.Balance})
}

func (s *BankHandler) HandleAcceptTransfer(c *gin.Context) {
	var acceptTransfer models.IncomingTransfer

	if err := c.ShouldBindJSON(&acceptTransfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, hasKey := c.Get("userID")
	if !hasKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userAccount, err := s.bankService.AcceptTransfer(acceptTransfer, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"New Balance": userAccount.Balance})
}
