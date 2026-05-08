package mq

import "context"

type Producer interface {
	PublishTask(ctx context.Context, message TaskMessage) error
}

type NoopProducer struct{}

func NewNoopProducer() Producer {
	return &NoopProducer{}
}

func (p *NoopProducer) PublishTask(ctx context.Context, message TaskMessage) error {
	// TODO: 替换为真实 MQ 投递逻辑。
	return nil
}
