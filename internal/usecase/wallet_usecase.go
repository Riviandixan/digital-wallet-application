package usecase

import (
	"context"
	"digital-wallet-application/internal/domain"

	"github.com/google/uuid"
)

type walletUsecase struct {
	walletRepo      domain.WalletRepository
	transactionRepo domain.TransactionRepository
}

func NewWalletUsecase(wr domain.WalletRepository, tr domain.TransactionRepository) domain.WalletUsecase {
	return &walletUsecase{
		walletRepo:      wr,
		transactionRepo: tr,
	}
}	

func (u *walletUsecase) GetBalance(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	return u.walletRepo.GetByUserID(ctx, userID)
}

func (u *walletUsecase) Withdraw(ctx context.Context, userID uuid.UUID, amount float64) (*domain.Wallet, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	return u.walletRepo.AtomicWithdraw(ctx, userID, amount)
}
