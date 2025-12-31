package processing

import "sync/atomic"

type Stats struct {
	processed atomic.Uint64
	failed    atomic.Uint64
}

func NewStats() *Stats { return &Stats{} }

func (s *Stats) IncProcessed() { s.processed.Add(1) }
func (s *Stats) IncFailed()    { s.failed.Add(1) }

func (s *Stats) Snapshot() (processed, failed uint64) {
	return s.processed.Load(), s.failed.Load()
}
