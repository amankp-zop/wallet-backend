package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/amankp-zop/wallet/internal/domain"
	"github.com/amankp-zop/wallet/internal/repository"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	ProcessTransferTask(ctx context.Context, t *asynq.Task) error
}

type RedisTaskProcessor struct {
	store repository.Store
}

func NewTaskProcessor(store repository.Store) TaskProcessor {
	return &RedisTaskProcessor{
		store: store,
	}
}

func (p *RedisTaskProcessor) ProcessTransferTask(ctx context.Context, t *asynq.Task) error{
	var payload ProcessTransferPayload
	if err:=json.Unmarshal(t.Payload(), &payload); err!= nil{
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	return p.store.ExecTx(ctx, func(q *repository.Queries) error {
		tx, err := q.GetTransactionByID(ctx, payload.TransactionID)
		if err != nil {
			return fmt.Errorf("failed to get transaction: %w", err)
		}

		if tx.Status != domain.TransactionStatusPending{
			log.Printf("Skipping already processed transaction %d with status %s", tx.ID, tx.Status)
			return nil
		}

		senderWallet, err := q.GetWalletForUpdate(ctx, tx.SenderWalletID)
		if err != nil {
			return fmt.Errorf("failed to get sender wallet: %w", err)
		}

		recieverWallet, err := q.GetWalletForUpdate(ctx, tx.ReceiverWalletID)
		if err != nil {
			return fmt.Errorf("failed to get receiver wallet: %w", err)
		}

		if senderWallet.Balance.LessThan(tx.Amount){
			err := q.UpdateTransactionStatus(ctx, tx.ID, domain.TransactionStatusFailed)
			if err != nil {
				return fmt.Errorf("failed to update transaction status: %w", err)
			}
	
			log.Printf("Transaction %d failed due to insufficient funds", tx.ID)
			return nil
		}

		err = q.UpdateWalletBalance(ctx, senderWallet.ID, senderWallet.Balance.Sub(tx.Amount))
		if err != nil {
			return fmt.Errorf("failed to update sender wallet balance: %w", err)
		}
		
		err = q.UpdateWalletBalance(ctx, recieverWallet.ID, recieverWallet.Balance.Add(tx.Amount))
		if err != nil {
			return fmt.Errorf("failed to update receiver wallet balance: %w", err)
		}

		return q.UpdateTransactionStatus(ctx, tx.ID, domain.TransactionStatusCompleted)
	})
}