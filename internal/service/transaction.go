package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/syarifid/bankx/internal/cfg"
	"github.com/syarifid/bankx/internal/dto"
	"github.com/syarifid/bankx/internal/ierr"
	"github.com/syarifid/bankx/internal/repo"
)

type TransactionService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newTransactionService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *TransactionService {
	return &TransactionService{repo, validator, cfg}
}

func (u *TransactionService) AddBalance(ctx context.Context, body dto.ReqAddBalance, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	bankAccountID := u.repo.BankAccount.GetBankAccountIDByNumber(ctx, body.SenderBankAccountNumber, body.SenderBankName)
	if bankAccountID == 0 {
		return ierr.ErrNotFound
	}
	transaction := body.ToTransactionEntity(sub, bankAccountID)
	err = u.repo.Transaction.AddBalance(ctx, sub, transaction)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}
