package config

import "os"

const (
	KafkaTopic = "statistics"
)

var (
	ServiceURL  = os.Getenv("SERVICE_URL")
	KafkaBroker = os.Getenv("KAFKA_BROKER")
	RedisAddress = os.Getenv("REDIS_URL")
)

