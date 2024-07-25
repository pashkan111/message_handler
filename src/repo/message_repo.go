package repo

import (
	"context"

	"messange_handler/src/entities"
	"messange_handler/src/errors/repo_errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message *entities.CreateMessageRequest,
) (int, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, repo_errors.OperationError{}
	}
	defer conn.Release()

	var message_id int
	err = conn.QueryRow(
		ctx,
		`INSERT INTO message (content)
		VALUES ($1)
		RETURNING message_id;
		`,
		message.Content,
	).Scan(&message_id)
	if err != nil {
		log.Error("Error with creating message:", err)
		return 0, repo_errors.OperationError{}
	}
	return message_id, nil
}
