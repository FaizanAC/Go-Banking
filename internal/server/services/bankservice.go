package services

import (
	"fmt"
	"time"

	"github.com/FaizanAC/Go-Banking/internal/models"
	"github.com/FaizanAC/Go-Banking/internal/util"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

const (
	DEPOSIT  = "DEPOSIT"
	WITHDRAW = "WITHDRAW"
	TRANSFER = "TRANSFER"
)

const (
	PENDING  = "PENDING"
	ACCEPTED = "ACCEPTED"
	EXPIRED  = "EXPIRED"
)

type BankService struct {
	db *gorm.DB
}

func NewBankService(db *gorm.DB) *BankService {
	return &BankService{db}
}

func (s *BankService) CreateAccount(userID uint) (models.BankAccount, error) {
	newAccount := models.BankAccount{
		AccountNumber: util.GenerateAccountNumber(),
		UserID:        userID,
		Balance:       0,
	}

	if res := s.db.Create(&newAccount); res.Error != nil {
		return newAccount, fmt.Errorf("failed to create account")
	}

	return newAccount, nil
}

func (s *BankService) GetAccountsByUserID(userID uint) ([]models.BankAccount, error) {
	var allAccounts []models.BankAccount

	if err := s.db.Where("user_id = ?", userID).Find(&allAccounts).Error; err != nil {
		return allAccounts, fmt.Errorf("user has no accounts")
	}

	return allAccounts, nil
}

func (s *BankService) DepositToAccount(deposit models.Transaction, userID uint) (models.BankAccount, error) {
	var account models.BankAccount
	if err := s.db.Where("account_number = ?", deposit.AccountNumber).First(&account).Error; err != nil {
		return account, fmt.Errorf("account not found")
	}

	if account.UserID != userID {
		return account, fmt.Errorf("you are not the owner of this account")
	}

	account.Balance += deposit.Amount

	deposit.Type = DEPOSIT
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
		return account, fmt.Errorf("failed to save deposit")
	}

	return account, nil
}

func (s *BankService) WithdrawFromAccount(withdraw models.Transaction, userID uint) (models.BankAccount, error) {
	var account models.BankAccount
	if err := s.db.Where("account_number = ?", withdraw.AccountNumber).First(&account).Error; err != nil {
		return account, fmt.Errorf("account not found")
	}

	if account.UserID != userID {
		return account, fmt.Errorf("you are not the owner of this account")
	}

	if account.Balance < withdraw.Amount {
		return account, fmt.Errorf("insufficient balance")
	}

	account.Balance -= withdraw.Amount

	withdraw.Type = WITHDRAW
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
		return account, fmt.Errorf("failed to save withdraw")
	}

	return account, nil
}

func (s *BankService) GetActivityFeed(userID uint) ([]models.Transaction, error) {
	var allAccounts []models.BankAccount
	if err := s.db.Where("user_id = ?", userID).Find(&allAccounts).Error; err != nil {
		return nil, fmt.Errorf("no accounts found")
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
		return nil, fmt.Errorf("failed to get activity feed")
	}

	return latestTransactions, nil
}

func (s *BankService) SendTransfer(transfer models.OutgoingTransfer, userID uint) (models.BankAccount, error) {
	var senderAccount models.BankAccount
	if err := s.db.Where("account_number = ?", transfer.AccountNumber).First(&senderAccount).Error; err != nil {
		return senderAccount, fmt.Errorf("sender account not found")
	}

	if senderAccount.UserID != userID {
		return senderAccount, fmt.Errorf("you are not the owner of this account")
	}

	transactionDetails := models.Transaction{
		Amount:        transfer.Amount,
		AccountNumber: transfer.AccountNumber,
		TransactionID: uuid.New().String(),
		Type:          TRANSFER,
	}

	transferRow := models.Transfer{
		SenderID:      userID,
		ReceiverID:    transfer.ReceiverID,
		Amount:        transfer.Amount,
		Status:        PENDING,
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
		return senderAccount, fmt.Errorf("failed to send transfer")
	}

	return senderAccount, nil
}

func (s *BankService) AcceptTransfer(acceptTransfer models.IncomingTransfer, userID uint) (models.BankAccount, error) {
	var tranferDetails models.Transfer
	if err := s.db.Where("transaction_id = ?", acceptTransfer.TransactionID).First(&tranferDetails).Error; err != nil {
		return models.BankAccount{}, fmt.Errorf("no transfer found")
	}

	if tranferDetails.ReceiverID != userID {
		return models.BankAccount{}, fmt.Errorf("you are not the receiver of this transfer")
	}

	var userAccount models.BankAccount
	if err := s.db.Where("account_number = ?", acceptTransfer.AccountNumber).Find(&userAccount).Error; err != nil {
		return userAccount, fmt.Errorf("account not found")
	}

	userAccount.Balance += tranferDetails.Amount
	tranferDetails.Status = ACCEPTED

	transactionDetails := models.Transaction{
		Amount:        tranferDetails.Amount,
		AccountNumber: userAccount.AccountNumber,
		TransactionID: uuid.New().String(),
		Type:          TRANSFER,
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
		return userAccount, fmt.Errorf("failed to accept transfer")
	}

	return userAccount, nil
}
