package pg

import (
	"fmt"

	"context"

	"messange_handler/config"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetPostgresPool(ctx context.Context, log *logrus.Logger) *pgxpool.Pool {
	postgresUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresDatabase,
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
