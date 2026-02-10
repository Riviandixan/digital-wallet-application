package repository

import (
	"context"
	"database/sql"
	"digital-wallet-application/internal/domain"
	"time"

	"github.com/google/uuid"
)

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) domain.WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	query := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE user_id = $1`

	var wallet domain.Wallet
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrWalletNotFound
	}

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepository) GetByUserIDWithLock(ctx context.Context, tx *sql.Tx, userID uuid.UUID) (*domain.Wallet, error) {
	query := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE user_id = $1 FOR UPDATE`

	var wallet domain.Wallet
	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, userID).Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.Balance,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
		)
	} else {
		err = r.db.QueryRowContext(ctx, query, userID).Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.Balance,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, domain.ErrWalletNotFound
	}

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, walletID uuid.UUID, newBalance float64) error {
	query := `UPDATE wallets SET balance = $1, updated_at = now() WHERE id = $2`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, newBalance, walletID)
	} else {
		_, err = r.db.ExecContext(ctx, query, newBalance, walletID)
	}
	return err
}

func (r *walletRepository) AtomicWithdraw(ctx context.Context, userID uuid.UUID, amount float64) (*domain.Wallet, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	wallet, err := r.GetByUserIDWithLock(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if wallet.Balance < amount {
		return nil, domain.ErrInsufficientFunds
	}

	newBalance := wallet.Balance - amount
	err = r.UpdateBalance(ctx, tx, wallet.ID, newBalance)
	if err != nil {
		return nil, err
	}

	txRecord := &domain.Transaction{
		ID:          uuid.New(),
		WalletID:    wallet.ID,
		Amount:      amount,
		Type:        domain.TransactionTypeWithdraw,
		Status:      domain.TransactionStatusSuccess,
		ReferenceID: uuid.New().String(),
		CreatedAt:   time.Now(),
	}

	queryLog := `INSERT INTO transactions (id, wallet_id, amount, type, status, reference_id, created_at) 
				 VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(ctx, queryLog,
		txRecord.ID,
		txRecord.WalletID,
		txRecord.Amount,
		txRecord.Type,
		txRecord.Status,
		txRecord.ReferenceID,
		txRecord.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	wallet.Balance = newBalance
	return wallet, nil
}
