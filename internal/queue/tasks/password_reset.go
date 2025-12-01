package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TaskSendResetPasswordEmail = "task:send_reset_password_email"
)


type ResetPasswordEmailPayload struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Code      string  `json:"code"`
}


func (q *Queue) EnqueueSendResetPasswordEmail(email, firstName, code string, delay time.Duration) error{
	payload, err := json.Marshal(ResetPasswordEmailPayload{FirstName: firstName , Code: code, Email: email })
    if err != nil{
		log.Errorw("failed to marshal payload", "error", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendResetPasswordEmail, payload)
    
	var opts []asynq.Option
	if delay > 0 {
	opts = append(opts, asynq.ProcessIn(delay))
}

info, err := q.client.Enqueue(task, opts...)
if err != nil{
	log.Errorw("failed to enquueu task", "error", err)
	return fmt.Errorf("failed to enqueue task: %w", err)
}

log.Infof("enqueue task: %s with ID: %s", info.Queue, info.ID)
return nil
}