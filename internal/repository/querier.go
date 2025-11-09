package repository

import "github.com/amankp-zop/wallet/internal/domain"

type Queries struct {
	domain.WalletRepository
	domain.UserRepository
	domain.TransactionRepository
}

func NewQueries(db DBTX) *Queries {
	return &Queries{
		WalletRepository: NewWalletRepository(db),
		UserRepository:   NewUserRepository(db),
		TransactionRepository: NewTransactionRepository(db),
	}
}
