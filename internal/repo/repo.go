package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	User  *userRepo
	Cat   *catRepo
	Match *matchRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.User = newUserRepo(conn)
	repo.Cat = newCatRepo(conn)
	repo.Match = newMatchRepo(conn)

	return &repo
}
