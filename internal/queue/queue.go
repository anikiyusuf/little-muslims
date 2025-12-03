package queue

import (
	"os"
	"os/signal"
	"syscall"

	
	"github.com/yusufaniki/muslim_tech/internal/queue/tasks"
	"github.com/yusufaniki/muslim_tech/internal/queue/workers"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)



func StartQueue(redisAddr, user, pass string,  worker *workers.Worker, logger *zap.SugaredLogger){
srv := asynq.NewServer(
	asynq.RedisClientOpt{
	Addr: redisAddr,
	Username: user,
	Password: pass,
},
	asynq.Config{
	Concurrency: 10, 

},
)



mux := asynq.NewServeMux()
mux.HandleFunc(tasks.TaskSendVerificationEmail,  worker.HandleSendVerificationEmail)
mux.HandleFunc(tasks.TaskSendResetPasswordEmail, worker.HandleSendResetPasswordEmail)

startTaskProcessing(srv, mux, logger)
}

func  startTaskProcessing(srv *asynq.Server, mux *asynq.ServeMux, logger *zap.SugaredLogger){
	if err := srv.Run(mux); err != nil{
		logger.Errorf("could not start queue: %v", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	srv.Shutdown()
}
