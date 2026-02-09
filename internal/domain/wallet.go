package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Wallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WalletUseCase interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	Withdraw(ctx context.Context, userID uuid.UUID, amount float64) (*Wallet, error)
}

type WalletRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	GetByUserIDWithLock(ctx context.Context, tx *sql.Tx, userID uuid.UUID) (*Wallet, error)
	UpdateBalance(ctx context.Context, tx *sql.Tx, walletID uuid.UUID, newBalance float64) error
	AtomicWithdraw(ctx context.Context, userID uuid.UUID, amoount float64) (*Wallet, error)
}
