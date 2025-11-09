package repository

import (
	"context"

	"github.com/amankp-zop/wallet/internal/domain"
)

type mysqlTransactionRepository struct {
	db DBTX
}

func NewTransactionRepository(db DBTX) domain.TransactionRepository {
	return &mysqlTransactionRepository{
		db: db,
	}
}

func (r *mysqlTransactionRepository)CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (sender_wallet_id, receiver_wallet_id, amount, status)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query, tx.SenderWalletID, tx.ReceiverWalletID, tx.Amount, tx.Status)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	tx.ID = id
	
	return nil
}