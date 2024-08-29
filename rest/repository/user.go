package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	postgres *database.PostgresDatabase
}

func NewUserRepository(postgres *database.PostgresDatabase) *UserRepository {
	return &UserRepository{
		postgres: postgres,
	}
}

func (repository *UserRepository) GetAllUsers() ([]*model.User, error) {
	rows, err := repository.postgres.Conn.Query(context.TODO(), selectAllUsers)
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

func (repository *UserRepository) GetUserById(id string) (*model.User, error) {
	row := repository.postgres.Conn.QueryRow(context.TODO(), selectUserById, id)
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

func (repository *UserRepository) Save(postUser *model.PostUser) (*model.User, error) {
	id := uuid.New().String()
	_, err := repository.postgres.Conn.Exec(context.TODO(), insertUser, id, postUser.Email)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    id,
		Email: postUser.Email,
	}, nil
}

func (repository *UserRepository) Update(id string, user *model.PostUser) (*model.User, error) {
	_, err := repository.postgres.Conn.Exec(context.TODO(), updateUser, user.Email, id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    id,
		Email: user.Email,
	}, nil
}

func (repository *UserRepository) Exists(id string) (bool, error) {
	_, err := repository.GetUserById(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repository *UserRepository) Delete(id string) error {
	_, err := repository.postgres.Conn.Exec(context.TODO(), deleteUser, id)
	if err != nil {
		return err
	}
	return nil
}
