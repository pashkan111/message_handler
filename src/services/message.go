package services

import (
	"context"
	"encoding/json"
	"messange_handler/config"
	"messange_handler/src/dependencies/kafka"
	"messange_handler/src/entities"
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
		kafka.ProduceMessage(config.MessageTopic, message_to_produce_json)
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

func UpdateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message_data []byte,
) error {
	var message entities.MessageToProduce
	err := json.Unmarshal(message_data, &message)
	if err != nil {
		log.Errorf("Error with unmarshalling message: %s", err)
		return err
	}
	err = repo.UpdateMessage(
		ctx,
		pool,
		log,
		message.MessageId,
	)
	if err != nil {
		log.Errorf("Error with updating message: %s", err)
		return err
	}
	return nil
}
