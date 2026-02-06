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
		WithQueueType("internal").
		WithMaxRetries(5).
		WithRetention(workqueue.NewDuration("168h")).
		WithScheduleIn(workqueue.NewDuration("10s")).
		Build()

	payload := workqueue.NewInputBuilder().
		WithOptions(opts).
		// WithEvent("payment.processed").
		WithEvent("event2").
		WithData(map[string]any{
			"input":    "value",
			"whoami":   "event1",
			"schedule": "10s",
		}).
		Build()

	return s.producer.Publish(ctx, payload)
}

func (s Service) Example02(ctx context.Context) error {
	payload := workqueue.NewInputBuilder().
		// WithEvent("payment.processed").
		WithEvent("event2").
		WithData(map[string]any{
			"input":  "value",
			"whoami": "event2",
		}).
		Build()

	return s.producer.Publish(ctx, payload)
}

func main() {
	ctx := context.Background()
	opts := workqueue.NewOptsBuilder().
		WithQueueType("external").
		WithMaxRetries(5).
		WithRetention(workqueue.NewDuration("168h")).
		WithScheduleIn(workqueue.NewDuration("5min")).
		Build()

	producer := SDK.NewProducer("my-app", "http://localhost:8080", "YWRtaW46cGFzc3dvcmQ=", opts)

	service := NewService(producer)

	if err := service.Example01(ctx); err != nil {
		panic(err)
	}

	if err := service.Example02(ctx); err != nil {
		panic(err)
	}
}
