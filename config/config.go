package config

import (
	"os"

	"github.com/joho/godotenv"
)

// kafka
var MessageTopic string
var MessageTopicGroupId string

var DeadQueueTopic string

var KafkaPort string
var KafkaHost string

// pg
var PostgresUser string
var PostgresPassword string
var PostgresHost string
var PostgresPort string
var PostgresDatabase string

func init() {
	godotenv.Load()

	MessageTopic = os.Getenv("KAFKA_MESSAGE_TOPIC_NAME")
	MessageTopicGroupId = os.Getenv("KAFKA_MESSAGE_TOPIC_GROUP_ID")

	DeadQueueTopic = os.Getenv("KAFKA_DEAD_QUEUE_TOPIC_NAME")

	KafkaPort = os.Getenv("KAFKA_PORT")
	KafkaHost = os.Getenv("KAFKA_HOST")

	PostgresUser = os.Getenv("PG_USER")
	PostgresPassword = os.Getenv("PG_PASSWORD")
	PostgresHost = os.Getenv("PG_HOST")
	PostgresPort = os.Getenv("PG_PORT")
	PostgresDatabase = os.Getenv("PG_DATABASE")
}
