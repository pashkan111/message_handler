package queue

import (
	"context"
	"fmt"
	"os"

	"messange_handler/src/entities"
	"messange_handler/src/repo"

	"github.com/jackc/pgx/v4/pgxpool"

	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MessageConsumer struct {
	Consumer *kafka.Consumer
	Log      *logrus.Logger
}

func (mc *MessageConsumer) InitConsumer() {
	bootstrap_servers := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	topic := os.Getenv("KAFKA_MESSAGE_TOPIC_NAME")
	group_id := os.Getenv("KAFKA_MESSAGE_TOPIC_GROUP_ID")

	var err error
	mc.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     bootstrap_servers,
		"auto.offset.reset":     "latest",
		"group.id":              group_id,
		"fetch.min.bytes":       1,
		"session.timeout.ms":    6000,
		"heartbeat.interval.ms": 2000,
	})
	if err != nil {
		mc.Log.Fatalf("Failed to create consumer: %s", err)
	}

	err = mc.Consumer.Subscribe(topic, nil)
	if err != nil {
		mc.Log.Fatalf("Failed to subscribe to topic: %s", err)
	}
}

func (mc *MessageConsumer) Run(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
) {
	mc.InitConsumer()
	mc.Log.Info("Message Consumer is running")
	go func() {
		defer mc.Consumer.Close()
		for {
			msg, err := mc.Consumer.ReadMessage(-1)
			if err != nil {
				mc.Log.Fatalf("consumer error: %v (%v)\n", err, msg)
			}

			// go func() {
			var message entities.MessageToProduce
			err = json.Unmarshal(msg.Value, &message)
			if err != nil {
				mc.Log.Errorf("Error with unmarshalling message: %s", err)
			} else {
				repo.UpdateMessage(
					ctx,
					pool,
					log,
					message.MessageId,
				)

			}
			// }()
		}
	}()
}
