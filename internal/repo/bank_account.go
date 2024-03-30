package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/syarifid/bankx/internal/dto"
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

func (br *bankAccountRepo) GetBalance(ctx context.Context, sub string) ([]dto.ResGetBalance, error) {
	q := `SELECT balance, currency FROM bank_accounts WHERE user_id = $1 ORDER BY balance DESC`

	rows, err := br.conn.Query(ctx, q, sub)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []dto.ResGetBalance
	for rows.Next() {
		var balance dto.ResGetBalance
		err = rows.Scan(&balance.Balance, &balance.Currency)
		if err != nil {
			return nil, err
		}

		balances = append(balances, balance)
	}

	return balances, nil
}
