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
	ReqTransaction struct {
		RecipientBankAccountNumber string `json:"recipientBankAccountNumber" validate:"required,min=5,max=30"`
		RecipientBankName          string `json:"recipientBankName" validate:"required,min=5,max=30"`
		FromCurrency               string `json:"fromCurrency" validate:"required,iso4217"`
		Balances                   string `json:"balances" validate:"required"`
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

func (d *ReqAddBalance) ToTransactionEntity(userId string, bankAccountId int) entity.Transaction {
	return entity.Transaction{
		Balance:          d.AddedBalance,
		Currency:         d.Currency,
		TransferProofImg: d.TransferProofImg,
		UserID:           userId,
		BankAccountID:    bankAccountId,
	}
}
