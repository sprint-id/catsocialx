package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/entity"
	"github.com/sprint-id/catsocialx/internal/ierr"
	timepkg "github.com/sprint-id/catsocialx/pkg/time"
)

type matchRepo struct {
	conn *pgxpool.Pool
}

func newMatchRepo(conn *pgxpool.Pool) *matchRepo {
	return &matchRepo{conn}
}

func (mr *matchRepo) MatchCat(ctx context.Context, sub string, match_cat entity.MatchCat) error {
	q := `INSERT INTO match_cats (user_id, match_cat_id, user_cat_id, message, created_at)
	VALUES ( $1, $2, $3, $4, EXTRACT(EPOCH FROM now())::bigint)`

	// show the query
	// fmt.Println(q)

	_, err := mr.conn.Exec(ctx, q,
		sub, match_cat.MatchCatId, match_cat.UserCatId, match_cat.Message)

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

func (mr *matchRepo) GetMatch(ctx context.Context, sub string) ([]dto.ResGetMatchCat, error) {
	q := `SELECT mc.id, u.name, u.email, u.created_at, 
		mcd.id, mcd.name, mcd.race, mcd.sex, mcd.description, mcd.age_in_month, mcd.image_urls, mcd.has_matched, mcd.created_at, 
		ucd.id, ucd.name, ucd.race, ucd.sex, ucd.description, ucd.age_in_month, ucd.image_urls, ucd.has_matched, ucd.created_at,
		mc.message, mc.created_at
		FROM match_cats mc
		INNER JOIN users u ON mc.user_id = u.id
		INNER JOIN cats mcd ON mc.user_cat_id = mcd.id
		INNER JOIN cats ucd ON mc.user_cat_id = ucd.id
		WHERE 1=1
		ORDER BY mc.created_at DESC`

	rows, err := mr.conn.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []dto.ResGetMatchCat
	for rows.Next() {
		var match dto.ResGetMatchCat
		var issuedByCreatedAt int64
		var matchCatCreatedAt int64
		var userCatCreatedAt int64
		var matchCreatedAt int64
		err = rows.Scan(&match.ID, &match.IssuedBy.Name, &match.IssuedBy.Email, &issuedByCreatedAt,
			&match.MatchCatDetail.ID, &match.MatchCatDetail.Name, &match.MatchCatDetail.Race, &match.MatchCatDetail.Sex, &match.MatchCatDetail.Description, &match.MatchCatDetail.AgeInMonth, &match.MatchCatDetail.ImageUrls, &match.MatchCatDetail.HasMatched, &matchCatCreatedAt,
			&match.UserCatDetail.ID, &match.UserCatDetail.Name, &match.UserCatDetail.Race, &match.UserCatDetail.Sex, &match.UserCatDetail.Description, &match.UserCatDetail.AgeInMonth, &match.UserCatDetail.ImageUrls, &match.UserCatDetail.HasMatched, &userCatCreatedAt,
			&match.Message, &matchCreatedAt)
		if err != nil {
			return nil, err
		}

		match.IssuedBy.CreatedAt = timepkg.TimeToISO8601(time.Unix(issuedByCreatedAt, 0))
		match.MatchCatDetail.CreatedAt = timepkg.TimeToISO8601(time.Unix(matchCatCreatedAt, 0))
		match.UserCatDetail.CreatedAt = timepkg.TimeToISO8601(time.Unix(userCatCreatedAt, 0))
		match.CreatedAt = timepkg.TimeToISO8601(time.Unix(matchCreatedAt, 0))
		matches = append(matches, match)
	}

	return matches, nil
}
