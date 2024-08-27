package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-examples/rest/model"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type InMemoryUserRepository struct {
	users map[string]*model.User
	mutex sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

func (r *InMemoryUserRepository) GetAllUsers() ([]*model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var users []*model.User
	for _, user := range r.users {
		users = append(users, user)
	}
	return nil, fmt.Errorf("io error")
}

func (r *InMemoryUserRepository) GetUserById(id string) (*model.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) Save(user *model.User) (*model.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, u := range r.users {
		if u.Email == user.Email {
			return nil, ErrUserAlreadyExists
		}
	}
	id := uuid.New().String()
	user.ID = id
	r.users[id] = user
	return user, nil
}

func (r *InMemoryUserRepository) Update(user *model.User) (*model.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.users[user.ID] = user
	return user, nil
}

func (r *InMemoryUserRepository) Exists(id string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	_, ok := r.users[id]
	return ok
}

func (r *InMemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.users, id)
	return nil
}
