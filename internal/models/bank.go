package models

import (
	"time"

	"gorm.io/gorm"
)

type BankAccount struct {
	gorm.Model
	AccountNumber string  `json:"accountNumber" binding:"required" gorm:"unique"`
	UserID        uint    `json:"userId" binding:"required"`
	Balance       float64 `json:"balance" binding:"required"`
}

type Transaction struct {
	gorm.Model
	Amount        float64   `json:"amount" binding:"required"`
	SenderID      uint      `json:"senderId" binding:"required"`
	ReceiverID    uint      `json:"receiverId" binding:"required"`
	TransactionID string    `json:"transactionId" gorm:"unique"`
	Status        string    `json:"status"`
	ExpiresOn     time.Time `json:"expiresOn"`
}
