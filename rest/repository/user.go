package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go-examples/rest/config"
	"go-examples/rest/database"
	"go-examples/rest/model"
)

var (
	selectAllUsers = "SELECT * FROM public.user"
	selectUserById = "SELECT * FROM public.user WHERE id = $1"
	insertUser     = "INSERT INTO public.user (id, email) VALUES ($1, $2)"
	updateUser     = "UPDATE public.user SET email = $1 WHERE id = $2"
	deleteUser     = "DELETE FROM public.user WHERE id = $1"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	database database.Database
	config   *config.DBConfig
}

func NewUserRepository(database database.Database, config *config.DBConfig) *UserRepository {
	return &UserRepository{
		database: database,
		config:   config,
	}
}

func (repository *UserRepository) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, repository.config.Timeout)
	defer cancel()
	rows, err := repository.database.Query(timeoutCtx, selectAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		user := new(model.User)
		err := rows.Scan(&user.ID, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repository *UserRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, repository.config.Timeout)
	defer cancel()
	row := repository.database.QueryRow(timeoutCtx, selectUserById, id)
	user := new(model.User)
	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (repository *UserRepository) Save(ctx context.Context, postUser *model.PostUser) (*model.User, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, repository.config.Timeout)
	defer cancel()
	id := uuid.New().String()
	_, err := repository.database.Exec(timeoutCtx, insertUser, id, postUser.Email)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    id,
		Email: postUser.Email,
	}, nil
}

func (repository *UserRepository) Update(ctx context.Context, id string, user *model.PostUser) (*model.User, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, repository.config.Timeout)
	defer cancel()
	_, err := repository.database.Exec(timeoutCtx, updateUser, user.Email, id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    id,
		Email: user.Email,
	}, nil
}

func (repository *UserRepository) Exists(ctx context.Context, id string) (bool, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, repository.config.Timeout)
	defer cancel()
	_, err := repository.GetUserById(timeoutCtx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repository *UserRepository) Delete(ctx context.Context, id string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, repository.config.Timeout)
	defer cancel()
	_, err := repository.database.Exec(timeoutCtx, deleteUser, id)
	if err != nil {
		return err
	}
	return nil
}
