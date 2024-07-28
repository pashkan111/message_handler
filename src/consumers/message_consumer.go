package consumers

import (
	"context"
	"fmt"

	"messange_handler/src/services"

	producer "messange_handler/src/dependencies/kafka"

	"github.com/jackc/pgx/v4/pgxpool"

	"messange_handler/config"

	"github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MessageConsumer struct {
	Consumer *kafka.Consumer
	Log      *logrus.Logger
}

func (mc *MessageConsumer) InitConsumer() {
	bootstrap_servers := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)

	var err error
	mc.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     bootstrap_servers,
		"auto.offset.reset":     "earliest",
		"group.id":              config.MessageTopicGroupId,
		"fetch.min.bytes":       1,
		"session.timeout.ms":    6000,
		"heartbeat.interval.ms": 2000,
	})
	if err != nil {
		mc.Log.Fatalf("Failed to create consumer: %s", err)
	}

	err = mc.Consumer.Subscribe(config.MessageTopic, nil)
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
				mc.Log.Errorf("error reading message: %v (%v)\n", err, msg)
			}

			go func(msg *kafka.Message) {
				err = services.UpdateMessage(ctx, pool, log, msg.Value)
				if err != nil {
					mc.Log.Errorf("error updating message: %v (%v)\n", err)
					producer.ProduceMessage(config.DeadQueueTopic, msg.Value)
				}

			}(msg)
		}
	}()
}
