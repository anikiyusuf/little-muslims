package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SQLRepository struct {
	connPool *pgxpool.Pool
	Queries  *Queries
}

var (
	dbInstance *pgxpool.Pool
)

func NewRepository(DBSource string, log *zap.SugaredLogger) (*SQLRepository, error) {
	if dbInstance != nil {
		return &SQLRepository{
			connPool: dbInstance,
			Queries:  New(dbInstance),
		}, nil
	}

    pgxpool, err := pgxpool.New(context.Background(), DBSource)
	if err != nil {
		log.Errorf("cannot connect to db : %v", err)
		return nil, err
	}


	dbInstance = pgxpool

	return &SQLRepository{
		connPool: pgxpool,
		Queries:  New(pgxpool),
	}, nil
}


func (s *SQLRepository) GetConnPool() *pgxpool.Pool {
	return s.connPool
}

