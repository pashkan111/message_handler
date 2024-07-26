package services

import (
	"context"
	"encoding/json"
	"messange_handler/src/entities"
	"messange_handler/src/queue"
	"messange_handler/src/repo"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message *entities.MessageFromRequest,
) (int, error) {
	message_id, err := repo.CreateMessage(ctx, pool, log, message)
	if err != nil {
		return 0, err
	}

	message_to_produce := entities.MessageToProduce{
		Content:   message.Content,
		MessageId: message_id,
	}
	message_to_produce_json, _ := json.Marshal(message_to_produce)

	go func() {
		queue.ProduceMessage(message_to_produce_json)
		log.Infof("Message sent to queue. Message id %d", message_id)
	}()
	return message_id, nil
}

func GetMessageStat(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
) (*entities.GetMessagesStatResponse, error) {
	stats, err := repo.GetMessageStats(ctx, pool, log)
	if err != nil {
		return nil, err
	}
	stats.TotalCount = stats.UnProcessedCount + stats.ProcessedCount
	return stats, nil
}
