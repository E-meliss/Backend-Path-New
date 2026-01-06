package eventstore

import (
	"context"
	"time"
)

type Event struct {
	Stream   string
	Type     string
	EntityID string
	Data     []byte
}

type AppendOptions struct {
	OccurredAt time.Time
}

type Store interface {
	Append(ctx context.Context, ev Event, opts AppendOptions) (int64, error)
	Load(ctx context.Context, stream string, fromID int64, limit int) ([]Event, error)
}
