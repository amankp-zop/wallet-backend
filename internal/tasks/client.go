package tasks

import "github.com/hibiken/asynq"

type TaskProducer interface {
	ProduceProcessTransferTask(transactionID int64) error
}

type RedisTaskProducer struct {
	client *asynq.Client
}

func NewTaskProducer(redisOpt asynq.RedisConnOpt) TaskProducer {
	client := asynq.NewClient(redisOpt);
	return &RedisTaskProducer{
		client: client,
	}
}

func (p *RedisTaskProducer) ProduceProcessTransferTask(transactinID int64) error {
	task, err:= NewProcessTransferTask(transactinID)
	if err != nil {
		return err
	}
	
	_,err = p.client.Enqueue(task)
	return err
}