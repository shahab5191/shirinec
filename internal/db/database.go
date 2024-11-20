package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/config"
)

type Database struct {
    Pool *pgxpool.Pool
}

func NewDatabase() (*Database, error) {
    pgxConfig, err := pgxpool.ParseConfig(config.AppConfig.DatabaseURL)
    if err != nil {
        return nil, err
    }

    pgxConfig.MaxConns = int32(config.AppConfig.PoolSize)
    pgxConfig.ConnConfig.ConnectTimeout = config.AppConfig.Timeout

    pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
    if err != nil {
        return nil, err
    }

    err = pool.Ping(context.Background())
    if err != nil {
        pool.Close()
        return nil, err
    }

    log.Println("Database connection established")
    return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
    db.Pool.Close()
}
