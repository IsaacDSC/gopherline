package main

import (
	"context"
	"github.com/IsaacDSC/workqueue"
	"github.com/IsaacDSC/workqueue/SDK"
)

type Service struct {
	producer workqueue.Producer
}

func NewService(producer workqueue.Producer) *Service {
	return &Service{producer: producer}
}

func (s Service) Example01(ctx context.Context) error {
	opts := workqueue.NewOptsBuilder().
		WithQueueType("internal.critical").
		WithMaxRetries(5).
		WithRetention(workqueue.NewDuration("168h")).
		WithScheduleIn(workqueue.NewDuration("5min")).
		Build()

	payload := workqueue.NewInputBuilder().
		WithOptions(opts).
		WithEvent("user.created").
		WithData(map[string]any{"input": "value"}).
		Build()

	return s.producer.Publish(ctx, payload)
}

func (s Service) Example02(ctx context.Context) error {
	payload := workqueue.NewInputBuilder().
		WithEvent("user.created").
		WithData(map[string]any{"input": "value"}).
		Build()

	return s.producer.Publish(ctx, payload)
}

func main() {
	ctx := context.Background()
	opts := workqueue.NewOptsBuilder().
		WithQueueType("internal.medium").
		WithMaxRetries(5).
		WithRetention(workqueue.NewDuration("168h")).
		WithScheduleIn(workqueue.NewDuration("5min")).
		Build()

	producer := SDK.NewProducer("http://localhost:8080", "your-token", opts)

	service := NewService(producer)

	if err := service.Example01(ctx); err != nil {
		panic(err)
	}

	if err := service.Example02(ctx); err != nil {
		panic(err)
	}
}
