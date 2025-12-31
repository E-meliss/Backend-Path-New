package domain

import (
	"errors"
	"sync"
	"time"
)

type Balance struct {
	UserID        int64 `json:"userId"`
	amount        Money
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	mu            sync.RWMutex
}

func NewBalance(userID int64, initial Money) *Balance {
	return &Balance{UserID: userID, amount: initial, LastUpdatedAt: time.Now()}
}

func (b *Balance) Get() Money {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.amount
}

func (b *Balance) Credit(x Money) error {
	if !x.IsPositive() {
		return errors.New("credit must be positive")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.amount += x
	b.LastUpdatedAt = time.Now()
	return nil
}

func (b *Balance) Debit(x Money) error {
	if !x.IsPositive() {
		return errors.New("debit must be positive")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.amount < x {
		return errors.New("insufficient funds")
	}
	b.amount -= x
	b.LastUpdatedAt = time.Now()
	return nil
}
