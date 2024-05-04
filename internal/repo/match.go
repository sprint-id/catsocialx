package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/catsocialx/internal/entity"
	"github.com/sprint-id/catsocialx/internal/ierr"
)

type matchRepo struct {
	conn *pgxpool.Pool
}

func newMatchRepo(conn *pgxpool.Pool) *matchRepo {
	return &matchRepo{conn}
}

func (mr *matchRepo) MatchCat(ctx context.Context, sub string, match_cat entity.MatchCat) error {
	q := `INSERT INTO match_cats (user_id, user_cat_id, message, created_at)
	VALUES ( $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint)`
	_, err := mr.conn.Exec(ctx, q,
		sub, match_cat.UserCatId, match_cat.Message)

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
