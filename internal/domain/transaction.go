package domain

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type TransactionType string

const (
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeDeposit  TransactionType = "deposit"
)

type TransactionStatus string

const (
	TransactionStatusSuccess TransactionStatus = "status"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	ID          uuid.UUID         `json:"id"`
	WalletID    uuid.UUID         `json:"wallet_id"`
	Amount      float64           `json:"amount"`
	Type        TransactionType   `json:"type"`
	Status      TransactionStatus `json:"status"`
	ReferenceID string            `json:"reference_id"`
	CreatedAt   time.Time         `json:"created_at"`
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
}
