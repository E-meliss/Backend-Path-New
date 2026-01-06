package eventstore

import (
	"context"
	"sync"
	"time"
)

type memoryStore struct {
	mu     sync.Mutex
	events []Event
}

func NewMemory() Store {
	return &memoryStore{}
}

func (m *memoryStore) Append(ctx context.Context, ev Event, opts AppendOptions) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if opts.OccurredAt.IsZero() {
		opts.OccurredAt = time.Now().UTC()
	}

	m.events = append(m.events, ev)
	return int64(len(m.events)), nil
}

func (m *memoryStore) Load(ctx context.Context, stream string, fromID int64, limit int) ([]Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var out []Event
	for _, e := range m.events {
		if e.Stream == stream {
			out = append(out, e)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}
