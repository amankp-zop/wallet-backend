package domain

import (
	"context"

	"github.com/shopspring/decimal"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	ID               int64             `json:"id"`
	SenderWalletID   int64             `json:"sender_wallet_id"`
	ReceiverWalletID int64             `json:"receiver_wallet_id"`
	Amount           decimal.Decimal   `json:"amount"`
	Status           TransactionStatus `json:"status"`
	CreatedAt        string            `json:"created_at"`
	UpdatedAt        string            `json:"updated_at"`
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *Transaction) error
}

type TransactionService interface {
	CreateTransfer(ctx context.Context, senderUserID, recieverUserID int64, amount decimal.Decimal) (*Transaction, error)
}