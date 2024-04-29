package dto

type (
	ReqTransaction struct {
		RecipientBankAccountNumber string `json:"recipientBankAccountNumber" validate:"required,min=5,max=30"`
		RecipientBankName          string `json:"recipientBankName" validate:"required,min=5,max=30"`
		FromCurrency               string `json:"fromCurrency" validate:"required,iso4217"`
		Balances                   int    `json:"balances" validate:"required"`
	}
)
