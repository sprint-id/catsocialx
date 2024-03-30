package service

import (
	"github.com/go-playground/validator/v10"

	"github.com/syarifid/bankx/internal/cfg"
	"github.com/syarifid/bankx/internal/repo"
)

type Service struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg

	User        *UserService
	Friend      *FriendService
	Post        *PostService
	Transaction *TransactionService
}

func NewService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *Service {
	service := Service{}
	service.repo = repo
	service.validator = validator
	service.cfg = cfg

	service.User = newUserService(repo, validator, cfg)
	service.Friend = newFriendService(repo, validator, cfg)
	service.Post = newPostService(repo, validator, cfg)
	service.Transaction = newTransactionService(repo, validator, cfg)

	return &service
}
