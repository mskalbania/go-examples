package test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go-examples/rest/model"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	args := u.Called()
	return args.Get(0).([]*model.User), args.Error(1)
}

func (u *UserRepositoryMock) GetUserById(ctx context.Context, id string) (*model.User, error) {
	args := u.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *UserRepositoryMock) Save(ctx context.Context, user *model.PostUser) (*model.User, error) {
	args := u.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *UserRepositoryMock) Update(ctx context.Context, id string, user *model.PostUser) (*model.User, error) {
	args := u.Called(id, user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *UserRepositoryMock) Exists(ctx context.Context, id string) (bool, error) {
	args := u.Called(id)
	return args.Bool(0), args.Error(1)
}

func (u *UserRepositoryMock) Delete(ctx context.Context, id string) error {
	args := u.Called(id)
	return args.Error(0)
}
