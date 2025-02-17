package models

import (
	"time"

	"gorm.io/gorm"
)

type BankAccount struct {
	gorm.Model
	AccountNumber string  `json:"accountNumber" gorm:"unique"`
	UserID        uint    `json:"userId"`
	Balance       float64 `json:"balance"`
}

type Transaction struct {
	gorm.Model
	Amount        float64 `json:"amount" binding:"required"`
	AccountNumber string  `json:"accountNumber" binding:"required"`
	TransactionID string  `json:"transactionId" gorm:"unique"`
	Type          string  `json:"type"`
}

type Transfer struct {
	gorm.Model
	SenderID      uint      `json:"senderId" binding:"required"`
	ReceiverID    uint      `json:"receiverId" binding:"required"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	ExpiresOn     time.Time `json:"expiresOn"`
	TransactionID string    `gorm:"unique"`
}
