package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 暴露给业务层使用的接口：有一个消息队列，底层是什么不重要，有两个功能
type Producer interface {
	PublishTask(ctx context.Context, message TaskMessage) error
	Close() error
}

type RabbitMQProducer struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	exchange   string
	routingKey string
	mu         sync.Mutex
}

func NewRabbitMQProducer(url, exchange, queue, routingKey string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	if err := ch.QueueBind(queue, routingKey, exchange, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("bind queue: %w", err)
	}

	return &RabbitMQProducer{
		conn:       conn,
		ch:         ch,
		exchange:   exchange,
		routingKey: routingKey,
	}, nil
}

func (p *RabbitMQProducer) PublishTask(ctx context.Context, message TaskMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal task message: %w", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.ch == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	if err := p.ch.PublishWithContext(ctx, p.exchange, p.routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	}); err != nil {
		return fmt.Errorf("publish task message: %w", err)
	}

	return nil
}

func (p *RabbitMQProducer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var closeErr error

	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			closeErr = fmt.Errorf("close rabbitmq channel: %w", err)
		}
		p.ch = nil
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil && closeErr == nil {
			closeErr = fmt.Errorf("close rabbitmq connection: %w", err)
		}
		p.conn = nil
	}

	return closeErr
}
