package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/yusufaniki/muslim_tech/internal/mailer"
	"github.com/yusufaniki/muslim_tech/internal/queue/tasks"
)

func (w *Worker) HandleSendVerificationEmail(ctx context.Context, t *asynq.Task) error{
var payload tasks.VerificationEmailPayload
if err := json.Unmarshal(t.Payload(), &payload); err != nil {
	w.logger.Errorw("could not unmarshal payload", "error", err)
	return fmt.Errorf("could not unmarshal payload: %w", err)
}
w.logger.Infof("sending verification email to %s with code %s", payload.Email)
data := map[string]interface{}{
	"FirstName":payload.FirstName,
	"Code":payload.Code,
	"Year": time.Now().Year(),
}

_, err := w.mailer.Send(mailer.VerificationEmail, payload.Email, data, true)

if err != nil {
	w.logger.Errorw("failed to send email", "email, payload.email")
	return fmt.Errorf("failed to send email: %w", err)
}
return nil 
}