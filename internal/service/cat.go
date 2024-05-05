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

func (u *CatService) AddCat(ctx context.Context, body dto.ReqAddOrUpdateCat, sub string) (dto.ResAddCat, error) {
	var res dto.ResAddCat
	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	cat := body.ToCatEntity(sub)
	res, err = u.repo.Cat.AddCat(ctx, sub, cat)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	return res, nil
}

func (u *CatService) GetCat(ctx context.Context, param dto.ParamGetCat, sub string) ([]dto.ResGetCat, error) {

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

func (u *CatService) GetCatByID(ctx context.Context, id, sub string) (dto.ResGetCat, error) {
	res, err := u.repo.Cat.GetCatByID(ctx, id, sub)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u *CatService) UpdateCat(ctx context.Context, body dto.ReqAddOrUpdateCat, id, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	cat := body.ToCatEntity(sub)
	err = u.repo.Cat.UpdateCat(ctx, id, sub, cat)
	if err != nil {
		return err
	}

	return nil
}

func (u *CatService) DeleteCat(ctx context.Context, id string, sub string) error {
	err := u.repo.Cat.DeleteCat(ctx, id, sub)
	if err != nil {
		return err
	}

	return nil
}
