package domain

import "context"

type Outbox struct {
	ID        int64  `json:"id"`
	Topic     string `json:"topic"`
	Payload   []byte `json:"payload"`
	Status    string `json:"status"`
}

type OutboxRepository interface {
	CreateOutbox(ctx context.Context, event *Outbox) error
	FetchAndLock(ctx context.Context, limit int64) ([]*Outbox, error)
	UpdateStatus(ctx context.Context, ids []int64, status string) error
}


