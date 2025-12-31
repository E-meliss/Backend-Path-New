package domain

import (
	"errors"
	"time"
)

type TransactionType string
type TransactionStatus string

const (
	TxCredit   TransactionType = "credit"
	TxDebit    TransactionType = "debit"
	TxTransfer TransactionType = "transfer"

	TxPending   TransactionStatus = "pending"
	TxCompleted TransactionStatus = "completed"
	TxFailed    TransactionStatus = "failed"
	TxReversed  TransactionStatus = "reversed"
)

type Transaction struct {
	ID         int64             `json:"id"`
	FromUserID *int64            `json:"fromUserId,omitempty"`
	ToUserID   *int64            `json:"toUserId,omitempty"`
	Amount     Money             `json:"amount"`
	Type       TransactionType   `json:"type"`
	Status     TransactionStatus `json:"status"`
	CreatedAt  time.Time         `json:"createdAt"`
}

func (t *Transaction) Validate() error {
	if !t.Amount.IsPositive() {
		return errors.New("amount must be positive")
	}
	switch t.Type {
	case TxCredit:
		if t.ToUserID == nil {
			return errors.New("credit requires to_user_id")
		}
	case TxDebit:
		if t.FromUserID == nil {
			return errors.New("debit requires from_user_id")
		}
	case TxTransfer:
		if t.FromUserID == nil || t.ToUserID == nil {
			return errors.New("transfer requires from_user_id and to_user_id")
		}
		if *t.FromUserID == *t.ToUserID {
			return errors.New("transfer requires different users")
		}
	default:
		return errors.New("invalid transaction type")
	}
	return nil
}

func (t *Transaction) MarkCompleted() error {
	if t.Status != TxPending {
		return errors.New("only pending tx can be completed")
	}
	t.Status = TxCompleted
	return nil
}

func (t *Transaction) MarkFailed() error {
	if t.Status != TxPending {
		return errors.New("only pending tx can be failed")
	}
	t.Status = TxFailed
	return nil
}

func (t *Transaction) MarkReversed() error {
	if t.Status != TxCompleted {
		return errors.New("only completed tx can be reversed")
	}
	t.Status = TxReversed
	return nil
}
