package repository

import (
	"context"
	"database/sql"

	"github.com/amankp-zop/wallet/internal/domain"
)

type walletRepository struct {
	db DBTX
}

func NewWalletRepository(db DBTX) domain.WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	query := `INSERT INTO wallets (user_id, balance, currency) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, wallet.UserID, wallet.Balance, wallet.Currency)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	wallet.ID = id

	return nil
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID int64) (*domain.Wallet, error) {
	query := `SELECT id, user_id, balance, currency, created_at, updated_at FROM wallets WHERE user_id = ?`
	row := r.db.QueryRowContext(ctx, query, userID)

	var wallet domain.Wallet
	err := row.Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &wallet, nil
}