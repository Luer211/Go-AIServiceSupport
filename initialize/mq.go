package initialize

import (
	"Go-AIServiceSupport/config"
	"Go-AIServiceSupport/internal/mq"
	"fmt"
)

func InitMQ(cfg *config.Config) (mq.Producer, error) {
	producer, err := mq.NewRabbitMQProducer(
		cfg.MQ.URL,
		cfg.MQ.Exchange,
		cfg.MQ.Queue,
		cfg.MQ.RoutingKey,
	)

	if err != nil {
		return nil, fmt.Errorf("initialize rabbitmq producer: %w", err)
	}

	return producer, nil
}
