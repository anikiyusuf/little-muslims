package tasks

import (
   	"github.com/hibiken/asynq"
)

type Queue struct {
	client *asynq.Client
}


func NewQueue(addr, user, pass string ) *Queue{
   return &Queue{client: asynq.NewClient(asynq.RedisClientOpt{Addr: addr, Username: user, Password: pass})}
}