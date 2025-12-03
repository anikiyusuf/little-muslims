package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/yusufaniki/muslim_tech/pkg/logger"
	"go.uber.org/zap"
)


const (
	TaskSendVerificationEmail = "task:send_verification_email"
)

var log *zap.SugaredLogger = logger.CreateZapLogger()

type VerificationEmailPayload struct {
	FirstName string  `json:"first_name"`
	Email     string  `json:"email"`
	Code      string  `json:"code"`
}

func (q *Queue) EnqueueSendVerificationEmail(email, firstName, code string, delay time.Duration) error {
	payload, err := json.Marshal(VerificationEmailPayload{FirstName: firstName, Code: code, Email: email})
	if err != nil {
		log.Errorw("failed to marshal payload","error", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerificationEmail, payload)

	var opts []asynq.Option
	if delay > 0 {
		opts = append(opts, asynq.ProcessIn(delay))
	}

	info, err := q.client.Enqueue(task, opts...)
	if err != nil {
		log.Errorw("failed to enqueue task", "error", err)
		return fmt.Errorf("failed to enqueue task :%w", err)	
	}

	log.Infof("enqueue task: %s with ID: %s", info.Queue,  info.ID)
	return nil
}