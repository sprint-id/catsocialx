package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/catsocialx/internal/cfg"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/ierr"
	"github.com/sprint-id/catsocialx/internal/repo"
)

type MatchService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newMatchService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *MatchService {
	return &MatchService{repo, validator, cfg}
}

func (u *MatchService) MatchCat(ctx context.Context, body dto.ReqMatchCat, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	// get name from sub
	name, err := u.repo.User.GetNameBySub(ctx, sub)
	if err != nil {
		return err
	}

	// get email from sub
	email, err := u.repo.User.GetEmailBySub(ctx, sub)
	if err != nil {
		return err
	}

	match := body.ToMatchCatEntity(name, email)
	err = u.repo.Match.MatchCat(ctx, sub, match)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}
