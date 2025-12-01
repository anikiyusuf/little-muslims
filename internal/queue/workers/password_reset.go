package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"


	"github.com/yusufaniki/muslim_tech/internal/queue/tasks"
	"github.com/yusufaniki/muslim_tech/internal/mailer"
	"github.com/hibiken/asynq"
)



func (w *Worker) HandleSendResetPasswordEmail(ctx context.Context, t *asynq.Task) error {
	var payload tasks.ResetPasswordEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		w.logger.Errorw("could not unmarshal payload", "error", err)
		return fmt.Errorf("could not unmarshal payload: %w", err)
	}
	w.logger.Infof("sending reset password email to %s with code %s", payload.Email, payload.Code)

	data := map[string]interface{}{
		"Firstname": payload.FirstName,
		"Code":      payload.Code,
		"Year":      time.Now().Year(),
	}

	_, err := w.mailer.Send(mailer.ResetPassword, payload.Email, data, true)
	
	if err != nil {
		w.logger.Errorw("failed to send email", "email", payload.Email, "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

return nil

}