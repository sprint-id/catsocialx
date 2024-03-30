package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type bankAccountRepo struct {
	conn *pgxpool.Pool
}

func newBankAccountRepo(conn *pgxpool.Pool) *bankAccountRepo {
	return &bankAccountRepo{conn}
}

func (br *bankAccountRepo) FindByUserIDAndCurrency(sub, currency string) (int, error) {
	q := `SELECT id FROM bank_accounts WHERE user_id = $1 AND currency = $2`

	var id int
	err := br.conn.QueryRow(context.Background(), q, sub, currency).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (br *bankAccountRepo) GetBankAccountIDByNumber(ctx context.Context, number, bankName string) int {
	q := `SELECT id FROM bank_accounts WHERE bank_account_number = $1 AND bank_name = $2`

	var id int
	err := br.conn.QueryRow(ctx, q, number, bankName).Scan(&id)
	if err != nil {
		return 0
	}

	return id
}
