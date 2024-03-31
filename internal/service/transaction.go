package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/syarifid/bankx/internal/cfg"
	"github.com/syarifid/bankx/internal/dto"
	"github.com/syarifid/bankx/internal/ierr"
	"github.com/syarifid/bankx/internal/repo"
	response "github.com/syarifid/bankx/pkg/resp"
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

	if body.TransferProofImg == "http://incomplete" {
		return ierr.ErrBadRequest
	}

	transaction := body.ToTransactionEntity(sub, body.SenderBankAccountNumber, body.SenderBankName)
	err = u.repo.Transaction.AddBalance(ctx, sub, transaction)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}

func (u *TransactionService) GetBalance(ctx context.Context, sub string) ([]dto.ResGetBalance, error) {
	balance, err := u.repo.Transaction.GetBalance(ctx, sub)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (u *TransactionService) GetBalanceHistory(ctx context.Context, param dto.ParamGetBalanceHistory, sub string) ([]dto.ResGetBalanceHistory, response.Meta, error) {
	meta := response.Meta{}

	if param.Limit == 0 {
		param.Limit = 5
	}

	err := u.validator.Struct(param)
	if err != nil {
		return nil, meta, ierr.ErrBadRequest
	}

	res, count, err := u.repo.Transaction.GetBalanceHistory(ctx, param, sub)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = count
	meta.Limit = param.Limit
	meta.Offset = param.Offset

	return res, meta, nil
}

func (u *TransactionService) AddTransaction(ctx context.Context, body dto.ReqTransaction, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	transaction := body.ToTransactionEntity(sub, body.RecipientBankAccountNumber, body.RecipientBankName)
	err = u.repo.Transaction.AddTransaction(ctx, sub, transaction)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}
