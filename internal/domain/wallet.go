package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *Wallet) error
	GetByUserID(ctx context.Context, userID int64) (*Wallet, error)
}

type WalletService interface {
	GetWalletByUserID(ctx context.Context, userID int64) (*Wallet, error)
}