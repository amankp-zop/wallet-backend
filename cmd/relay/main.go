package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/amankp-zop/wallet/internal/config"
	"github.com/amankp-zop/wallet/internal/database"
	"github.com/amankp-zop/wallet/internal/repository"
	"github.com/amankp-zop/wallet/internal/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	log.Println("Starting Relay Process...")

	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewDatabase(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connected Successfully.")
	defer db.Close()

	redisOpt, err := asynq.ParseRedisURI(cfg.Redis.Addr)
	if err != nil {
		log.Fatalf("Could not parse redis url: %v", err)
	}

	store := repository.NewStore(db)
	taskProducer := tasks.NewTaskProducer(redisOpt)


	relay := NewRelay(store, taskProducer)
	relay.Run(context.Background())
}

type Relay struct {
	store        repository.Store
	taskProducer tasks.TaskProducer
}

func NewRelay(store repository.Store, taskProducer tasks.TaskProducer) *Relay {
	return &Relay{
		store:        store,
		taskProducer: taskProducer,
	}
}

func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			err:= r.processBatch(ctx)
			if err != nil {
				log.Printf("Error processing batch: %v", err)
			}
		case <- ctx.Done():
			log.Println("Relay process shutting down...")
			return
		}
	}
}

func (r *Relay) processBatch(ctx context.Context) error {
	return r.store.ExecTx(ctx, func(q *repository.Queries) error {
		events, err:= q.FetchAndLock(ctx, 100);
		if err != nil {
			return fmt.Errorf("failed to fetch and lock events: %w", err)
		}

		if len(events)>0{
			log.Printf("Relay found %d events to publish", len(events))
		}

		var processedIDs []int64
		for _, event := range events{
			payload := tasks.ProcessTransferPayload{}
			if err:= json.Unmarshal(event.Payload, &payload); err!=nil{
				log.Printf("failed to unmarshal event payload: %v", err)
				continue
			}

			err := r.taskProducer.ProduceProcessTransferTask(payload.TransactionID)
			if err != nil {
				log.Printf("failed to produce process transfer task: %v", err)
				continue
			}

			processedIDs = append(processedIDs, event.ID)
		}

		if len(processedIDs)>0{
			err = q.UpdateStatus(ctx, processedIDs, "PUBLISHED")
			if err != nil{
				return fmt.Errorf("failed to update events status: %w", err)
			}
		}

		return nil
	})
}