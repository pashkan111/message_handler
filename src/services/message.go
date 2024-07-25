package services

import (
	"context"

	"messange_handler/src/entities"
	"messange_handler/src/repo"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message *entities.CreateMessageRequest,
) (int, error) {
	message_id, err := repo.CreateMessage(ctx, pool, log, message)
	if err != nil {
		return 0, err
	}
	return message_id, nil
}
