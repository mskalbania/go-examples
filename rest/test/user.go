package test

import (
	"github.com/stretchr/testify/mock"
	"go-examples/rest/model"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) GetAllUsers() ([]*model.User, error) {
	args := u.Called()
	return args.Get(0).([]*model.User), args.Error(1)
}

func (u *UserRepositoryMock) GetUserById(id string) (*model.User, error) {
	args := u.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *UserRepositoryMock) Save(user *model.PostUser) (*model.User, error) {
	args := u.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *UserRepositoryMock) Update(id string, user *model.PostUser) (*model.User, error) {
	args := u.Called(id, user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *UserRepositoryMock) Exists(id string) (bool, error) {
	args := u.Called(id)
	return args.Bool(0), args.Error(1)
}

func (u *UserRepositoryMock) Delete(id string) error {
	args := u.Called(id)
	return args.Error(0)
}
