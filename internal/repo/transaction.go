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
	q := `INSERT INTO transactions (user_id, bank_name, bank_account_number, balance, currency, transfer_proof_img, created_at)
	VALUES ( $1, $2, $3, $4, $5, $6, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	_, err := tr.conn.Exec(ctx, q,
		sub, transaction.Source.BankName, transaction.Source.BankAccountNumber, transaction.Balance, transaction.Currency, transaction.TransferProofImg)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (tr *transactionRepo) GetBalance(ctx context.Context, sub string) ([]dto.ResGetBalance, error) {
	q := `SELECT SUM(balance), currency FROM transactions WHERE user_id = $1 GROUP BY currency ORDER BY SUM(balance) DESC`

	rows, err := tr.conn.Query(ctx, q, sub)
	if err != nil {
		return nil, err
	}

	var res []dto.ResGetBalance
	for rows.Next() {
		var r dto.ResGetBalance
		err = rows.Scan(&r.Balance, &r.Currency)
		if err != nil {
			return nil, err
		}

		res = append(res, r)
	}

	return res, nil
}

func (tr *transactionRepo) GetBalanceHistory(ctx context.Context, param dto.ParamGetBalanceHistory, sub string) ([]dto.ResGetBalanceHistory, int, error) {
	q := `SELECT id, bank_account_number, bank_name, balance, currency, transfer_proof_img, created_at FROM transactions
		WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := tr.conn.Query(ctx, q, sub, param.Limit, param.Offset)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.ResGetBalanceHistory
	for rows.Next() {
		var r dto.ResGetBalanceHistory
		err = rows.Scan(&r.ID, &r.Source.BankAccountNumber, &r.Source.BankName, &r.Balance, &r.Currency, &r.TransferProofImg, &r.CreatedAt)
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
