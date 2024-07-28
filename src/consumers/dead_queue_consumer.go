package consumers

import (
	"context"
	"fmt"
	"time"

	"messange_handler/src/services"

	"github.com/jackc/pgx/v4/pgxpool"

	"messange_handler/config"

	"github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type DeadQueueConsumer struct {
	Consumer *kafka.Consumer
	Log      *logrus.Logger
}

func (dqc *DeadQueueConsumer) InitConsumer() {
	bootstrap_servers := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)

	var err error
	dqc.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     bootstrap_servers,
		"auto.offset.reset":     "earliest",
		"group.id":              config.MessageTopicGroupId,
		"fetch.min.bytes":       1,
		"session.timeout.ms":    6000,
		"heartbeat.interval.ms": 2000,
		"enable.auto.commit":    false,
	})
	if err != nil {
		dqc.Log.Fatalf("Failed to create consumer: %s", err)
	}

	err = dqc.Consumer.Subscribe(config.DeadQueueTopic, nil)
	if err != nil {
		dqc.Log.Fatalf("Failed to subscribe to topic: %s", err)
	}
}

func (dqc *DeadQueueConsumer) Run(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
) {
	dqc.InitConsumer()
	dqc.Log.Info("Dead queue Consumer is running")

	go func() {
		defer dqc.Consumer.Close()
		for {
			msg, err := dqc.Consumer.ReadMessage(-1)
			if err != nil {
				dqc.Log.Errorf("error reading message: %v (%v)\n", err, msg)
			}

			retries := 3
			success := false
			for i := 0; i < retries; i++ {
				err = services.UpdateMessage(ctx, pool, log, msg.Value)
				if err != nil {
					time.Sleep(2 * time.Second)
					continue
				}
				break
			}
			if !success {
				dqc.Log.Errorf("error processing message in dead queue: %v (%v)\n", err, msg)
			}

			dqc.Consumer.CommitMessage(msg)
		}
	}()
}
