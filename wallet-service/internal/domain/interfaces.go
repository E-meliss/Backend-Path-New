package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, u User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, t Transaction) (Transaction, error)
	UpdateStatus(ctx context.Context, id int64, status TransactionStatus) error
	GetByID(ctx context.Context, id int64) (Transaction, error)
}

type BalanceRepository interface {
	GetForUpdate(ctx context.Context, userID int64) (Money, error)
	Set(ctx context.Context, userID int64, amount Money) error
}

type AuditRepository interface {
	Log(ctx context.Context, entityType, entityID, action string, details any) error
}

type UserService interface {
	Register(ctx context.Context, username, email, password string) (User, error)
	Authenticate(ctx context.Context, email, password string) (User, error)
	Authorize(u User, required UserRole) error
}

type TransactionService interface {
	Credit(ctx context.Context, toUserID int64, amount Money) (Transaction, error)
	Debit(ctx context.Context, fromUserID int64, amount Money) (Transaction, error)
	Transfer(ctx context.Context, fromUserID, toUserID int64, amount Money) (Transaction, error)
	Rollback(ctx context.Context, txID int64) error
}

type BalanceService interface {
	Get(ctx context.Context, userID int64) (Money, error)
}
