package dbrepo

import (
	"github.com/gergab1129/bookings/internal/config"
	"github.com/gergab1129/bookings/internal/repository"
	"github.com/jmoiron/sqlx"
)

type postgresDBrepo struct {
	App *config.AppConfig
	DB  *sqlx.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sqlx.DB
}

func NewPostgresRepo(conn *sqlx.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBrepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
