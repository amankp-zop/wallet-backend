package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TaskTypeProcessTransfer = "transfer:process"

type ProcessTransferPayload struct {
	TransactionID int64 `json:"transaction_id"`
}

func NewProcessTransferTask(transcationID int64) (*asynq.Task, error){
	payload ,err:= json.Marshal(ProcessTransferPayload{
		TransactionID: transcationID,
	})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskTypeProcessTransfer, payload), nil
}