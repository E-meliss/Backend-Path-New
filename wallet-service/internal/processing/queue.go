package processing

import (
	"context"

	"github.com/E-meliss/wallet-service/internal/domain"
)

type TaskType string

const (
	TaskCredit   TaskType = "credit"
	TaskDebit    TaskType = "debit"
	TaskTransfer TaskType = "transfer"
)

type Task struct {
	Type       TaskType
	FromUserID int64
	ToUserID   int64
	Amount     domain.Money
	ResultCh   chan error
}

type Queue struct {
	ch chan Task
}

func NewQueue(buffer int) *Queue {
	return &Queue{ch: make(chan Task, buffer)}
}

func (q *Queue) Enqueue(ctx context.Context, t Task) error {
	select {
	case q.ch <- t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *Queue) Chan() <-chan Task { return q.ch }
func (q *Queue) Close()            { close(q.ch) }
