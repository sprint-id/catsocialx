package entity

type SourceTransaction struct {
	ID                int    `json:"id"`
	BankAccountName   string `json:"bank_account_name"`
	BankAccountNumber string `json:"bank_account_number"`

	TransactionID string `json:"transaction_id"`
}
