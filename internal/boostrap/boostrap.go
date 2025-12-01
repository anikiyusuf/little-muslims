package boostrap

import (
	"fmt"
	"time"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yusufaniki/muslim_tech/internal/cache"
	"github.com/yusufaniki/muslim_tech/internal/config"
	"github.com/yusufaniki/muslim_tech/internal/mailer"
	"github.com/yusufaniki/muslim_tech/internal/queue"
	"github.com/yusufaniki/muslim_tech/internal/queue/tasks"
	"github.com/yusufaniki/muslim_tech/internal/queue/workers"
	"github.com/yusufaniki/muslim_tech/internal/repository"
	"github.com/yusufaniki/muslim_tech/pkg/auth"
	"github.com/yusufaniki/muslim_tech/pkg/logger"
	"go.uber.org/zap"
)

type Application struct {
	Config      *config.Config
	ConnPool    *pgxpool.Pool
	Logger      *zap.SugaredLogger
	JWTManager  *auth.JWTManager
	Cache       *cache.RedisCache
	Queue		*tasks.Queue
	Mailer		*mailer.EmailClient
	Worker      *workers.Worker
	Repository  *repository.Queries

}


func InitializeApp() (*Application, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	//Create Logger instance
	appLogger := logger.CreateZapLogger()
	appLogger.Info("Configuration loaded successfully")


	// Database store
	repository, err := repository.NewRepository(cfg.DBSource, appLogger)
	if err != nil {
		appLogger.Errorf("Failed to initialize database repository: %v", err)
		return nil, fmt.Errorf("failed to initialize database repository: %v", err)
	}
    
	connPool:= repository.GetConnPool()
	appLogger.Info("Database store initialized successfully")

	var redisCache *cache.RedisCache
	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	if cfg.Redis.Enabled{
		redisCache = cache.NewRedisCache(redisAddr, cfg.Redis.Pass)
		appLogger .Info("Cache initialized successfully")
	}

newQueue := tasks.NewQueue(redisAddr, cfg.Redis.User, cfg.Redis.Pass)
appLogger.Info("Queue initialized successfully")


mailtrapProvider :=  mailer.NewMailtrapProvider(cfg.Mail.Host, cfg.Mail.User, cfg.Mail.Pass, cfg.Mail.From, cfg.Mail.Port)
emailClient := mailer.NewEmailClient(mailtrapProvider, appLogger)
appLogger.Info("Mailer initialized successfully")

   // Set up queue worker
   newWorker := workers.NewWorker(emailClient, appLogger)

go queue.StartQueue(redisAddr, cfg.Redis.User, cfg.Redis.Pass, newWorker, appLogger)
appLogger.Info("Queue worker started successfully")

jwtManager := auth.NewJWTManager(cfg.JWT.Secret, time.Duration(cfg.JWT.Duration)*time.Hour, *repository.Queries)
appLogger.Info("JWT Manager initialized successfully")



	return &Application{
		Config:     &cfg,
		Logger:     appLogger,
		ConnPool:   connPool,
		JWTManager: jwtManager,
		Cache:      redisCache,
		Queue:      newQueue,
		Mailer:     emailClient,
		Worker:		newWorker,
		Repository: repository.Queries,
	}, nil
}