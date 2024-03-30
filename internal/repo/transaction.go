package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syarifid/bankx/internal/dto"
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
	VALUES ( $1, $2, $3, $4, $5, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

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

func (tr *transactionRepo) GetBalanceHistory(ctx context.Context, param dto.ParamGetBalanceHistory, sub string) ([]dto.ResGetBalanceHistory, int, error) {
	q := `SELECT transactions.id, transactions.balance, transactions.currency, transactions.transfer_proof_img, transactions.created_at, bank_accounts.bank_account_number, bank_accounts.bank_name
		FROM transactions
		JOIN bank_accounts ON transactions.bank_account_id = bank_accounts.id
		WHERE transactions.user_id = $1
		ORDER BY transactions.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := tr.conn.Query(ctx, q, sub, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.ResGetBalanceHistory
	for rows.Next() {
		var r dto.ResGetBalanceHistory
		err = rows.Scan(&r.ID, &r.Balance, &r.Currency, &r.TransferProofImg, &r.CreatedAt, &r.Source.BankAccountNumber, &r.Source.BankName)
		if err != nil {
			return nil, 0, err
		}

		res = append(res, r)
	}

	q = `SELECT count(id) FROM transactions WHERE user_id = $1`
	var count int
	err = tr.conn.QueryRow(ctx, q, sub).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return res, count, nil
}
