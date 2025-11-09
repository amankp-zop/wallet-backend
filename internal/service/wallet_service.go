package service

import (
	"context"
	"errors"

	"github.com/amankp-zop/wallet/internal/domain"
	"github.com/amankp-zop/wallet/internal/repository"
)

var (
	ErrWalletNotFound = errors.New("wallet Not found")
)

type walletService struct {
	store   repository.Store
}

func NewWalletService(store repository.Store) domain.WalletService {
	return &walletService{
		store: store,
	}
}

func (s *walletService) GetWalletByUserID(ctx context.Context, userID int64) (*domain.Wallet, error) {
	wallet, err := s.store.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	return wallet, nil
}
