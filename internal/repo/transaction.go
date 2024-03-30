package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syarifid/bankx/internal/entity"
	"github.com/syarifid/bankx/internal/ierr"
)

type transactionRepo struct {
	conn *pgxpool.Pool
}

func newTransactionRepo(conn *pgxpool.Pool) *transactionRepo {
	return &transactionRepo{conn}
}

func (tr *transactionRepo) AddBalance(ctx context.Context, sub string, transaction entity.Transaction) error {
	// add transaction
	q := `INSERT INTO transactions (user_id, bank_account_id, balance, currency, transfer_proof_img, created_at)
	VALUES ( $1, $2, $3, $4, $5, now()) RETURNING id`

	_, err := tr.conn.Exec(ctx, q,
		sub, transaction.BankAccountID, transaction.Balance, transaction.Currency, transaction.TransferProofImg)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	// update balance
	q = `UPDATE bank_accounts SET balance = balance + $1 WHERE id = $2 AND currency = $3`

	_, err = tr.conn.Exec(ctx, q,
		transaction.Balance, transaction.BankAccountID, transaction.Currency)

	if err != nil {
		return err
	}

	return nil
}
