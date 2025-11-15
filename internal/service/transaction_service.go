package service

import (
	"context"
	"encoding/json"

	"github.com/amankp-zop/wallet/internal/domain"
	"github.com/amankp-zop/wallet/internal/repository"
	"github.com/amankp-zop/wallet/internal/tasks"
	"github.com/shopspring/decimal"
)

type transactionService struct {
	store repository.Store
}

func NewTransactionService(store repository.Store) domain.TransactionService {
	return &transactionService{
		store: store,
	}
}

func (s *transactionService) CreateTransfer(ctx context.Context, senderUserID, receiverUserID int64, amount decimal.Decimal) (*domain.Transaction, error){
	var createdTx *domain.Transaction

	err := s.store.ExecTx(ctx, func(q *repository.Queries) error {
		senderWallet, err := s.store.GetByUserID(ctx, senderUserID);
		if err != nil {
			return err
		}

		receiverWallet, err := s.store.GetByUserID(ctx, receiverUserID);
		if err != nil {
			return err
		}

		tx := &domain.Transaction{
			SenderWalletID: senderWallet.ID,
			ReceiverWalletID: receiverWallet.ID,
			Amount: amount,
			Status: domain.TransactionStatusPending,
		}

		err = s.store.CreateTransaction(ctx, tx)
		if err != nil {
			return err
		}

		payload := tasks.ProcessTransferPayload{
			TransactionID: tx.ID,
		}
		payloadBytes, err:= json.Marshal(payload);
		if err != nil {
			return err
		}

		outboxEvent := &domain.Outbox{
			Topic: tasks.TaskTypeProcessTransfer,
			Payload: payloadBytes,
		}

		if err := q.CreateOutbox(ctx, outboxEvent); err!=nil{
			return err
		}

		createdTx = tx
		return nil
	})

	return createdTx, err
}