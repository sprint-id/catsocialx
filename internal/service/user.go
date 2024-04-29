package service

import (
	"context"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/catsocialx/internal/cfg"
	"github.com/sprint-id/catsocialx/internal/dto"
	"github.com/sprint-id/catsocialx/internal/ierr"
	"github.com/sprint-id/catsocialx/internal/repo"
	"github.com/sprint-id/catsocialx/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newUserService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *UserService {
	return &UserService{repo, validator, cfg}
}

func (u *UserService) Register(ctx context.Context, body dto.ReqRegister) (dto.ResRegister, error) {
	res := dto.ResRegister{}

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	_, err = mail.ParseAddress(body.Email)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	user := body.ToEntity(u.cfg.BCryptSalt)
	userID, err := u.repo.User.Insert(ctx, user)
	if err != nil {
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: userID})
	if err != nil {
		return res, err
	}

	res.Email = body.Email
	res.Name = body.Name
	res.AccessToken = token

	return res, nil
}

func (u *UserService) Login(ctx context.Context, body dto.ReqLogin) (dto.ResLogin, error) {
	res := dto.ResLogin{}

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	_, err = mail.ParseAddress(body.Email)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	user, err := u.repo.User.GetByEmail(ctx, body.Email)
	if err != nil {
		return res, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return res, err
	}

	res.Email = user.Email
	res.Name = user.Name
	res.AccessToken = token

	return res, nil
}

func (u *UserService) UpdateAccount(ctx context.Context, body dto.ReqUpdateAccount, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	if body.ImageURL == "http://incomplete" {
		return ierr.ErrBadRequest
	}

	err = u.repo.User.LookUp(ctx, sub)
	if err != nil {
		return err
	}

	err = u.repo.User.UpdateAccount(ctx, sub, body.Name, body.ImageURL)
	return err
}
