package postgres

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func NewDB(ctx context.Context, config *DBConfig) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode)

	return pgxpool.New(ctx, connectionString)
}

func NewDBFromEnv(ctx context.Context) (*pgxpool.Pool, error) {
	return NewDB(ctx, &DBConfig{
		Host: os.Getenv("DB_HOST"),
		Port: func() int {
			port, err := strconv.Atoi(os.Getenv("DB_PORT"))
			if err != nil {
				return 0
			}
			return port
		}(),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	})
}
