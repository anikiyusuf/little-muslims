package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yusufaniki/muslim_tech/pkg/logger"
)



type HealthService struct {
	ConnPool *pgxpool.Pool
}


func NewHealthService(connPool *pgxpool.Pool) *HealthService {
	return &HealthService{ConnPool: connPool}
}

func (h *HealthService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := h.ConnPool.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] ="db is up and running"

	// Get Pool statistics
	dbStats := h.ConnPool.Stat()

	stats["total_connection"] = strconv.Itoa(int(dbStats.TotalConns()))
	stats["idle_connection"] = strconv.Itoa(int(dbStats.IdleConns()))
	stats["in_use_connections"] = strconv.Itoa(int(dbStats.AcquiredConns()))
	stats["max_connection"] = strconv.Itoa(int(dbStats.MaxConns()))

	// Additional statistics
	stats["new_connections"] = strconv.Itoa(int(dbStats.NewConnsCount()))
	stats["canceled_acquires"] = strconv.Itoa(int(dbStats.CanceledAcquireCount()))
	stats["acquire_count"] =  strconv.Itoa(int(dbStats.AcquireCount()))
	stats["acquire_duration"] = dbStats.AcquireDuration().String()
	stats["max_idle_closed"] = strconv.Itoa(int(dbStats.MaxIdleDestroyCount()))
	stats["max_lifetime_closed"] = strconv.Itoa(int(dbStats.MaxLifetimeDestroyCount()))
	stats["successful_acquires"]  = strconv.Itoa(int(dbStats.AcquireCount() - dbStats.CanceledAcquireCount()))

	// Add Custom messages based on thresholds
	if dbStats.TotalConns() > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.CanceledAcquireCount() > 100 {
		stats["message"] = "High number of canceled acquire attempts. Review connection pool configurations."
	}

	if dbStats.MaxIdleDestroyCount() > int64(dbStats.IdleConns()/2){
		stats["message"] = "Many idle connections are being closed. Consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeDestroyCount() > int64(dbStats.TotalConns()/2){
		stats["message"] = "Many connections are  being closed due to max lifetime. Consider increasing max lifetime or revising usage pattern."
	}

	return stats

}

func (h *HealthService) Close() error {
	logger.CreateZapLogger().Info("Disconnected from  database")
	h.ConnPool.Close()
	return nil
}