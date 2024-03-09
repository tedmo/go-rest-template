package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tedmo/go-rest-template/internal/app"
	"github.com/tedmo/go-rest-template/internal/postgres/sqlc"
)

type UserRepo struct {
	DB      *pgxpool.Pool
	Queries sqlc.Querier
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		DB:      db,
		Queries: sqlc.New(),
	}
}

func (repo *UserRepo) FindUserByID(ctx context.Context, id int64) (*app.User, error) {
	sqlUser, err := repo.Queries.FindUserByID(ctx, repo.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return sqlUser.DomainModel(), nil
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *app.CreateUserReq) (*app.User, error) {

	sqlUser, err := repo.Queries.CreateUser(ctx, repo.DB, user.Name)
	if err != nil {
		return nil, err
	}

	return sqlUser.DomainModel(), nil
}

func (repo *UserRepo) FindUsers(ctx context.Context) ([]app.User, error) {

	sqlUsers, err := repo.Queries.FindUsers(ctx, repo.DB)
	if err != nil {
		return nil, err
	}

	var users []app.User
	for _, sqlUser := range sqlUsers {
		users = append(users, *sqlUser.DomainModel())
	}

	if users == nil {
		users = []app.User{}
	}

	return users, nil
}
