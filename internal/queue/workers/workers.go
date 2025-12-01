package workers

import (
	"github.com/yusufaniki/muslim_tech/internal/mailer"
	"go.uber.org/zap"
)

type Worker struct {
	mailer *mailer.EmailClient
	logger  *zap.SugaredLogger
}




func NewWorker(mailClient *mailer.EmailClient, logger *zap.SugaredLogger) *Worker{
	return &Worker{mailer: mailClient, logger: logger}
}