package repository

import (
	"context"
	"database/sql"
	"digital-wallet-application/internal/domain"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (id, wallet_id, amount, type, status, reference_id, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		tx.ID,
		tx.WalletID,
		tx.Amount,
		tx.Type,
		tx.Status,
		tx.ReferenceID,
		tx.CreatedAt,
	)
	return err
}
