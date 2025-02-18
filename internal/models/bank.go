package models

import (
	"time"
)

type BankAccount struct {
	GormModel
	AccountNumber string  `json:"accountNumber" gorm:"unique"`
	UserID        uint    `json:"userId"`
	Balance       float64 `json:"balance"`
}

type Transaction struct {
	GormModel
	Amount        float64 `json:"amount" binding:"required"`
	AccountNumber string  `json:"accountNumber" binding:"required"`
	TransactionID string  `json:"transactionId" gorm:"unique"`
	Type          string  `json:"type"`
}

type Transfer struct {
	GormModel
	SenderID      uint      `json:"senderId" binding:"required"`
	ReceiverID    uint      `json:"receiverId" binding:"required"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	ExpiresOn     time.Time `json:"expiresOn"`
	TransactionID string    `gorm:"unique"`
}

type OutgoingTransfer struct {
	Amount        float64 `json:"amount" binding:"required"`
	AccountNumber string  `json:"accountNumber" binding:"required"`
	ReceiverID    uint    `json:"receiverID" binding:"required"`
}

type IncomingTransfer struct {
	TransactionID string `json:"transactionID" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
}
