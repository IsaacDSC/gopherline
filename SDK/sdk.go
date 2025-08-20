package SDK

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IsaacDSC/workqueue"
	"io"
	"net/http"
	"time"
)

type Producer struct {
	host   string
	token  string
	client *http.Client
	opts   workqueue.Opts
}

var _ workqueue.Producer = (*Producer)(nil)

func NewProducer(host string, token string, opts workqueue.Opts) *Producer {
	client := &http.Client{
		Timeout: time.Duration(50 * time.Millisecond),
	}
	return &Producer{host: host, token: token, client: client, opts: opts}
}

func (p Producer) Publish(ctx context.Context, input workqueue.Input) error {
	if input.Event == "" {
		return fmt.Errorf("event cannot be empty")
	}

	var emptyOpts workqueue.Opts
	if input.Options == emptyOpts {
		input.Options = p.opts
	}

	return p.publish(ctx, workqueue.Payload{
		Event:   input.Event,
		Data:    input.Data,
		Options: input.Options,
		Metadata: map[string]any{
			"headers": map[string]string{
				"correlation_id": input.CorrelationID,
				"event_id":       input.EventID,
			},
		},
	})
}

func (p Producer) publish(ctx context.Context, payload workqueue.Payload) error {
	url := fmt.Sprintf("%s/event/publisher", p.host)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", p.token))

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed on publisher: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to cofirmation ack: %w", err)
	}

	if resp.StatusCode > 399 {
		return fmt.Errorf("error on publisher: %s", string(body))
	}

	return nil
}
