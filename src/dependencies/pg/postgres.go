package pg

import (
	"fmt"
	"os"

	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetPostgresPool(ctx context.Context, log *logrus.Logger) *pgxpool.Pool {
	postgresUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DATABASE"),
	)
	config, err := pgxpool.ParseConfig(postgresUrl)
	if err != nil {
		log.Fatal("Error with parsing config", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Error with connecting to DB", err)
	}
	return pool
}
