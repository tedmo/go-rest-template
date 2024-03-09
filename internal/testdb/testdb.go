package testdb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultDBUsername     = "username"
	defaultDBPassword     = "password"
	defaultDBName         = "postgres"
	defaultStartupTimeout = 5 * time.Second
)

type TestDB struct {
	pool              *pgxpool.Pool
	postgresContainer *postgres.PostgresContainer
}

// New launches a new TestDB.  The caller should call Close when finished to shut it down.
func New(ctx context.Context) (*TestDB, error) {

	// Postgres
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase(defaultDBName),
		postgres.WithUsername(defaultDBUsername),
		postgres.WithPassword(defaultDBPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(defaultStartupTimeout)),
	)
	if err != nil {
		return nil, fmt.Errorf("error starting postgres container: %w", err)
	}

	return &TestDB{
		postgresContainer: postgresContainer,
	}, nil
}

// Open returns a connection pool for the test database instance.
// If the connection pool is not closed explicitly, it will be closed when the TestDB.Close method is called.
func (db *TestDB) Open(ctx context.Context) (*pgxpool.Pool, error) {
	if db.pool != nil {
		if err := db.pool.Ping(ctx); err != nil {
			return db.pool, nil
		}
	}

	connectionString, err := db.ConnectionString()
	if err != nil {
		return nil, err
	}

	db.pool, err = pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	return db.pool, nil
}

// ConnectionString returns the connection string for connecting to the test database instance
func (db *TestDB) ConnectionString() (string, error) {
	connectionString, err := db.postgresContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		return "", fmt.Errorf("error getting db connection string: %w", err)
	}
	return connectionString, nil
}

// Migrate runs DB migrations found in the migrationsPath file path
func (db *TestDB) Migrate(migrationsPath string) error {
	postgresURL, err := db.ConnectionString()
	if err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}

	m, err := migrate.New("file://"+migrationsPath, postgresURL)
	if err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}

	return m.Up()
}

// Close closes the test DB connection pool and shuts down the test DB instance.
func (db *TestDB) Close(ctx context.Context) error {
	if db.pool != nil {
		db.pool.Close()
	}

	if db.postgresContainer != nil && db.postgresContainer.IsRunning() {
		if err := db.postgresContainer.Terminate(ctx); err != nil {
			return fmt.Errorf("error closing db: %w", err)
		}
	}

	return nil
}
