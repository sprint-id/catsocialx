package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/catsocialx/internal/cfg"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/ierr"
	"github.com/sprint-id/catsocialx/internal/repo"
)

type CatService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newCatService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *CatService {
	return &CatService{repo, validator, cfg}
}

func (u *CatService) AddCat(ctx context.Context, body dto.ReqAddCat, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
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
	cat := body.ToCatEntity(sub)
	err = u.repo.Cat.AddCat(ctx, sub, cat)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}

func (u *CatService) GetCat(ctx context.Context, param dto.ParamGetCat, sub string) ([]dto.ResGetCat, error) {

	if param.Limit == 0 {
		param.Limit = 5
	}

	err := u.validator.Struct(param)
	if err != nil {
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Cat.GetCat(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
