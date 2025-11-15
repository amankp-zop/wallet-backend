// File: internal/repository/db.go

package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/amankp-zop/wallet/internal/domain" // Make sure domain is imported
)

// Store defines all the functions to execute db queries and transactions.
// IT NOW EMBEDS THE REPOSITORY INTERFACES.
type Store interface {
	ExecTx(ctx context.Context, fn func(*Queries) error) error
	domain.UserRepository   
	domain.WalletRepository
	domain.TransactionRepository
	domain.OutboxRepository
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: NewQueries(db),
	}
}

// ExecTx executes a function within a database transaction
func (s *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := NewQueries(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}