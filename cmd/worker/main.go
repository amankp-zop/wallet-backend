package main

import (
	"log"

	"github.com/hibiken/asynq"

	"github.com/amankp-zop/wallet/internal/config"
	"github.com/amankp-zop/wallet/internal/database"
	"github.com/amankp-zop/wallet/internal/repository"
	"github.com/amankp-zop/wallet/internal/tasks"
)

func main() {
	log.Println("Starting Worker Process...")

	cfg, err := config.LoadConfig("./configs")
	if err!=nil{
		log.Fatalf("Error loading config: %v", err)
	}

	db, err:= database.NewDatabase(cfg.Database.DSN)
	if err!= nil{
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connected Successfully.")
	defer db.Close()

	redisOpt, err := asynq.ParseRedisURI(cfg.Redis.Addr)
	if err != nil {
		log.Fatalf("Could not parse redis url: %v", err)
	}

	store:= repository.NewStore(db)
	taskProcessor:= tasks.NewTaskProcessor(store)

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,

			Queues: map[string]int{
				"default": 1,
			},
		},
	)

	mux:=asynq.NewServeMux()
	mux.HandleFunc(tasks.TaskTypeProcessTransfer, taskProcessor.ProcessTransferTask)

	log.Println("Worker is running...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run worker server: %v", err)
	}
}