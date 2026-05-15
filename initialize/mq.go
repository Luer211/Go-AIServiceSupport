package initialize

import (
	"Go-AIServiceSupport/config"
	"Go-AIServiceSupport/internal/mq"
)

func InitMQ(cfg *config.Config) mq.Producer {
	producer, err := mq.NewRabbitMQProducer(
		cfg.MQ.URL,
		cfg.MQ.Exchange,
		cfg.MQ.Queue,
		cfg.MQ.RoutingKey,
	)
	// Todo：这里应该要报错
	if err != nil {
		return nil
	}

	return producer
}
