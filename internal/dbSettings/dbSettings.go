package dbSettings

import (
	"DockerPostgreExample/internal/logger"
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func Initialize() (*pgxpool.Pool, error) {
	dir, err := os.Getwd()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}
	environmentPath := filepath.Join(dir, ".env")
	err = godotenv.Load(environmentPath) // load .env
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}

	// change POSTGRES_HOST=localhost if run not in docker container
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))

	dbpool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("Unable to connect to database")
		os.Exit(1)
	}
	logger.Log.Info().Msg("Database connection established")

	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()

	err = createTablesIfNotExists(conn) //create tables if not exists
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}

	return dbpool, err
}

//go:embed createTablesIfNotExists.sql
var createTablesSql string

func createTablesIfNotExists(db *pgxpool.Conn) error {
	_, err := db.Exec(context.Background(), createTablesSql)
	return err
}
