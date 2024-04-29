package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/entity"
	"github.com/sprint-id/catsocialx/internal/ierr"
	timepkg "github.com/sprint-id/catsocialx/pkg/time"
)

type catRepo struct {
	conn *pgxpool.Pool
}

func newCatRepo(conn *pgxpool.Pool) *catRepo {
	return &catRepo{conn}
}

// {
// 	"name": "", // not null, minLength 1, maxLength 30
// 	"race": "", /** not null, enum of:
// 			- "Persian"
// 			- "Maine Coon"
// 			- "Siamese"
// 			- "Ragdoll"
// 			- "Bengal"
// 			- "Sphynx"
// 			- "British Shorthair"
// 			- "Abyssinian"
// 			- "Scottish Fold"
// 			- "Birman" */
// 	"sex": "", // not null, enum of: "male" / "female"
// 	"ageInMonth": 1, // not null, min: 1, max: 120082
// 	"description":"" // not null, minLength 1, maxLength 200
// 	"imageUrls":[ // not null, minItems: 1, items: not null, should be url
// 		"","",""
// 	]
// }

func (cr *catRepo) AddCat(ctx context.Context, sub string, cat entity.Cat) error {
	// add cat
	q := `INSERT INTO cats (user_id, name, race, sex, age_in_month, description, image_urls, created_at)
	VALUES ( $1, $2, $3, $4, $5, $6, $7, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	image_urls := "{" + strings.Join(cat.ImageUrls, ",") + "}" // Format image URLs as a PostgreSQL array

	_, err := cr.conn.Exec(ctx, q,
		sub, cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, image_urls)

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

func (cr *catRepo) GetCat(ctx context.Context, param dto.ParamGetCat, sub string) ([]dto.ResGetCat, error) {
	var query strings.Builder

	if param.IsAlreadyMatched {
		query.WriteString("AND is_already_matched = true ")
	} else {
		query.WriteString("AND is_already_matched = false ")
	}

	if param.Owned {
		query.WriteString("AND user_id = $1 ")
	} else {
		query.WriteString("AND user_id != $1 ") // apakah termasuk yang owned kalau misalkan owned = false?
	}

	if param.Search != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Search)))
	}

	rows, err := cr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.ResGetCat, 0, 10)
	for rows.Next() {
		var imageUrl sql.NullString
		var createdAt time.Time

		result := dto.ResGetCat{}
		err := rows.Scan(&result.ID, &result.Name, &imageUrl, &createdAt, &result.IsAlreadyMatched)
		if err != nil {
			return nil, err
		}

		result.ImageUrls = strings.Split(imageUrl.String, ",")
		result.CreatedAt = timepkg.TimeToISO8601(createdAt)
		results = append(results, result)
	}

	return results, nil
}
