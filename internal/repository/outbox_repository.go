package repository

import (
	"context"

	"github.com/amankp-zop/wallet/internal/domain"
)

type mysqlOutboxRepository struct {
	db DBTX
}

func NewOutboxRepository(db DBTX) domain.OutboxRepository {
	return &mysqlOutboxRepository{
		db: db,
	}
}

func (r *mysqlOutboxRepository) CreateOutbox(ctx context.Context, event *domain.Outbox) error {
	query := "INSERT INTO outbox (topic, payload) VALUES (?, ?)"
	_, err := r.db.ExecContext(ctx, query, event.Topic, event.Payload)
	
	return err
}