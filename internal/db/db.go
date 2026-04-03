package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var Pool *pgxpool.Pool

func ConnectDB(connStr string, logger *zap.Logger) {
	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		logger.Fatal("Database Ping failed", zap.Error(err))
	}

	logger.Info("Successfully connected to database")
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}

func GetPool() *pgxpool.Pool {
	return Pool
}
