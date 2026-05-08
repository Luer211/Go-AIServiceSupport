package initialize

import (
	"Go-AIServiceSupport/config"
	"Go-AIServiceSupport/internal/mq"
)

func InitMQ(cfg *config.Config) mq.Producer {
	// TODO: 替换成真实 MQ producer，例如 RabbitMQ/Kafka/Redis Stream。
	return mq.NewNoopProducer()
}
