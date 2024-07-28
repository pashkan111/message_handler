package kafka

import (
	"fmt"
	"messange_handler/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

var producer *kafka.Producer

func InitProducer(log *logrus.Logger) {
	bootstrap_servers := fmt.Sprintf("%s:%s", config.KafkaHost, config.KafkaPort)

	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrap_servers,
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					log.Infof("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()
}

func ProduceMessage(topic string, message []byte) {
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}, nil)
}
