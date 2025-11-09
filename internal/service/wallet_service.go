package service

import (
	"context"
	"errors"

	"github.com/amankp-zop/wallet/internal/domain"
)

var (
	ErrWalletNotFound = errors.New("wallet Not found")
)

type walletService struct {
	walletRepo domain.WalletRepository
}

func NewWalletService(repo domain.WalletRepository) domain.WalletService {
	return &walletService{
		walletRepo: repo,
	}
}

func (s *walletService) GetWalletByUserID(ctx context.Context, userID int64) (*domain.Wallet, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	return wallet, nil
}
