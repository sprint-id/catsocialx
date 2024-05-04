package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/entity"
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

// ResAddCat struct {
// 	ID        string `json:"id"`
// 	CreatedAt string `json:"createdAt"`
// }

func (cr *catRepo) AddCat(ctx context.Context, sub string, cat entity.Cat) (dto.ResAddCat, error) {
	// add cat
	q := `INSERT INTO cats (user_id, name, race, sex, age_in_month, description, image_urls, created_at)
	VALUES ( $1, $2, $3, $4, $5, $6, $7, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	image_urls := "{" + strings.Join(cat.ImageUrls, ",") + "}" // Format image URLs as a PostgreSQL array

	var id string
	err := cr.conn.QueryRow(ctx, q, sub,
		cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, image_urls).Scan(&id)
	if err != nil {
		return dto.ResAddCat{}, err
	}

	createdAt := time.Now()
	return dto.ResAddCat{ID: id, CreatedAt: timepkg.TimeToISO8601(createdAt)}, nil
}

func (cr *catRepo) GetCat(ctx context.Context, param dto.ParamGetCat, sub string) ([]dto.ResGetCat, error) {
	var query strings.Builder

	if param.HasMatched {
		query.WriteString(`SELECT c.id, c.name, c.image_urls, c.created_at, EXISTS (
			SELECT 1 FROM match_cats m WHERE m.user_cat_id = c.id AND m.user_id = $1
		) AS has_matched, c.description FROM cats c WHERE EXISTS (
			SELECT 1 FROM match_cats m WHERE m.user_cat_id = c.id
		)`)
	} else {
		query.WriteString(`SELECT c.id, c.name, c.image_urls, c.created_at, EXISTS (
			SELECT 1 FROM match_cats m WHERE m.user_cat_id = c.id AND m.user_id = $1
		) AS has_matched, c.description FROM cats c WHERE NOT EXISTS (
			SELECT 1 FROM match_cats m WHERE m.user_cat_id = c.id
		)`)
	}

	// param id
	if param.ID != "" {
		id, err := strconv.Atoi(param.ID)
		if err != nil {
			return nil, err
		}
		query.WriteString(fmt.Sprintf("AND id = %d", id))
	}

	if param.Owned {
		query.WriteString("AND user_id = $1 ")
	} else {
		query.WriteString("AND user_id != $1 ") // apakah termasuk yang owned kalau misalkan owned = false?
	}

	if param.Search != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Search)))
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	// show query
	// fmt.Println(query.String())

	rows, err := cr.conn.Query(ctx, query.String(), sub) // Replace $1 with sub
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.ResGetCat, 0, 10)
	for rows.Next() {
		var imageUrl sql.NullString
		var createdAt int64
		var description string

		result := dto.ResGetCat{}
		err := rows.Scan(&result.ID, &result.Name, &imageUrl, &createdAt, &result.HasMatched, &description)
		if err != nil {
			return nil, err
		}

		result.ImageUrls = strings.Split(imageUrl.String, ",")
		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		result.Description = description
		results = append(results, result)
	}

	return results, nil
}

func (cr *catRepo) GetCatByID(ctx context.Context, id string, sub string) (dto.ResGetCat, error) {
	q := `SELECT id,
		name,
		race,
		sex,
		age_in_month,
		description,
		image_urls,
		EXISTS (
			SELECT 1 FROM match_cats m WHERE m.user_cat_id = c.id AND m.user_id = $1
		) AS has_matched,
		created_at
	FROM cats c WHERE id = $2`

	var imageUrl sql.NullString
	var createdAt int64
	var description string

	result := dto.ResGetCat{}
	err := cr.conn.QueryRow(ctx, q, sub, id).Scan(&result.ID, &result.Name, &result.Race, &result.Sex, &result.AgeInMonth, &description, &imageUrl, &result.HasMatched, &createdAt)
	if err != nil {
		return dto.ResGetCat{}, err
	}

	result.ImageUrls = strings.Split(imageUrl.String, ",")
	result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
	result.Description = description

	return result, nil
}

func (cr *catRepo) UpdateCat(ctx context.Context, id, sub string, cat entity.Cat) error {
	q := `UPDATE cats SET 
		name = $1,
		race = $2,
		sex = $3,
		age_in_month = $4,
		description = $5,
		image_urls = $6
	WHERE
		id = $7 AND user_id = $8`

	image_urls := "{" + strings.Join(cat.ImageUrls, ",") + "}" // Format image URLs as a PostgreSQL array

	_, err := cr.conn.Exec(ctx, q,
		cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, image_urls, id, sub)

	if err != nil {
		return err
	}

	return nil
}

func (cr *catRepo) DeleteCat(ctx context.Context, id string, sub string) error {
	q := `DELETE FROM cats WHERE id = $1 AND user_id = $2`

	// log id and sub
	fmt.Println("id: ", id)
	fmt.Println("sub: ", sub)

	_, err := cr.conn.Exec(ctx, q, id, sub)
	if err != nil {
		return err
	}

	return nil
}
