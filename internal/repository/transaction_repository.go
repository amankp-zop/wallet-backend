package repository

import (
	"context"
	"database/sql"

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

func (r *mysqlTransactionRepository)GetTransactionByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	query := `
		SELECT id, sender_wallet_id, receiver_wallet_id, amount, status, created_at, updated_at
		FROM transactions
		WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var tx domain.Transaction
	err := row.Scan(&tx.ID, &tx.SenderWalletID, &tx.ReceiverWalletID, &tx.Amount, &tx.Status, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *mysqlTransactionRepository)UpdateTransactionStatus(ctx context.Context, id int64, status domain.TransactionStatus) error {
	query := `
		UPDATE transactions
		SET status = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}