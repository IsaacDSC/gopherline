package gopherline

import "context"

type Producer interface {
	Publish(ctx context.Context, payload Input) error
}
