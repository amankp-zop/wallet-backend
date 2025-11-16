package repository

import (
	"context"
	"fmt"
	"strings"

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

func (r *mysqlOutboxRepository) FetchAndLock(ctx context.Context, limit int64) ([]*domain.Outbox, error) {
	query := `SELECT id, topic, payload
		FROM outbox
		WHERE status= 'UNPUBLISHED'
		ORDER BY id
		LIMIT ?
		FOR UPDATE SKIP LOCKED`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*domain.Outbox
	for rows.Next(){
		var event domain.Outbox
		if err := rows.Scan(&event.ID, &event.Topic, &event.Payload); err!=nil{
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}

func (r *mysqlOutboxRepository) UpdateStatus(ctx context.Context, ids []int64, status string) error {
	if len(ids) ==0{
		return nil
	}

	query := fmt.Sprintf("UPDATE outbox SET status = ? WHERE id IN (%s)",
		strings.Repeat("?,", len(ids)-1)+"?",
	)

	args := make([]any, len(ids)+1)
	args[0] = status
	for i, id := range ids {
		args[i+1] = id
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	
	return err
}