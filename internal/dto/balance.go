package dto

import (
	"github.com/syarifid/bankx/internal/entity"
)

type (
	ReqAddBalance struct {
		SenderBankAccountNumber string `json:"senderBankAccountNumber" validate:"required,min=5,max=30"`
		SenderBankName          string `json:"senderBankName" validate:"required,min=5,max=30"`
		AddedBalance            int    `json:"addedBalance" validate:"required,min=1"`
		Currency                string `json:"currency" validate:"required,iso4217"`
		TransferProofImg        string `json:"transferProofImg" validate:"required,url"`
	}
	ParamGetBalanceHistory struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}
	ResGetBalance struct {
		Balance  int    `json:"balance"`
		Currency string `json:"currency"`
	}
	ResGetBalanceHistory struct {
		ID               string        `json:"id"`
		Balance          int           `json:"balance"`
		Currency         string        `json:"currency"`
		TransferProofImg string        `json:"transferProofImg"`
		CreatedAt        int64         `json:"createdAt"`
		Source           entity.Source `json:"source"`
	}
)

func (d *ReqAddBalance) ToTransactionEntity(userId, bankAccountNumber, bankName string) entity.Transaction {
	return entity.Transaction{
		Balance:          d.AddedBalance,
		Currency:         d.Currency,
		TransferProofImg: d.TransferProofImg,
		Source: entity.Source{
			BankAccountNumber: bankAccountNumber,
			BankName:          bankName,
		},
		UserID: userId,
	}
}

func (d *ReqTransaction) ToTransactionEntity(userId, bankAccountNumber, bankName string) entity.Transaction {
	return entity.Transaction{
		Balance:          d.Balances,
		Currency:         d.FromCurrency,
		TransferProofImg: "",
		Source: entity.Source{
			BankAccountNumber: bankAccountNumber,
			BankName:          bankName,
		},
		UserID: userId,
	}
}
