package main

import (
	"context"
	"github.com/IsaacDSC/gopherline"
	"github.com/IsaacDSC/gopherline/SDK"
)

type Service struct {
	producer gopherline.Producer
}

func NewService(producer gopherline.Producer) *Service {
	return &Service{producer: producer}
}

func (s Service) Example01(ctx context.Context) error {
	opts := gopherline.NewOptsBuilder().
		WithQueueType("internal.critical").
		WithMaxRetries(5).
		WithRetention(gopherline.NewDuration("168h")).
		WithScheduleIn(gopherline.NewDuration("5min")).
		Build()

	payload := gopherline.NewInputBuilder().
		WithOptions(opts).
		WithEvent("user.created").
		WithData(map[string]any{"input": "value"}).
		Build()

	return s.producer.Publish(ctx, payload)
}

func (s Service) Example02(ctx context.Context) error {
	payload := gopherline.NewInputBuilder().
		WithEvent("user.created").
		WithData(map[string]any{"input": "value"}).
		Build()

	return s.producer.Publish(ctx, payload)
}

func main() {
	ctx := context.Background()
	opts := gopherline.NewOptsBuilder().
		WithQueueType("internal.medium").
		WithMaxRetries(5).
		WithRetention(gopherline.NewDuration("168h")).
		WithScheduleIn(gopherline.NewDuration("5min")).
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
