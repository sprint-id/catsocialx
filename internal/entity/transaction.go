package entity

type Transaction struct {
	ID               string `json:"id"`
	Balance          int    `json:"balance"`
	Currency         string `json:"currency"`
	TransferProofImg string `json:"transfer_proof_img"`
	CreatedAt        int64  `json:"created_at"`
	Source           Source `json:"source"`

	UserID string `json:"user_id"`
}
