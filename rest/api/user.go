package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-examples/rest/model"
	"go-examples/rest/repository"
	"net/http"
)

type UserRepository interface {
	GetAllUsers() ([]*model.User, error)
	GetUserById(id string) (*model.User, error)
	Save(user *model.PostUser) (*model.User, error)
	Update(id string, user *model.PostUser) (*model.User, error)
	Exists(id string) (bool, error)
	Delete(id string) error
}

type UserAPI struct {
	userRepository UserRepository
}

func NewUserAPI(userRepository UserRepository) *UserAPI {
	return &UserAPI{userRepository: userRepository}
}

func (userAPI *UserAPI) GetUsers(context *gin.Context) {
	users, err := userAPI.userRepository.GetAllUsers()
	if err != nil {
		AbortWithContextError(context, http.StatusInternalServerError, "error getting users", err)
		return
	}
	context.JSON(http.StatusOK, users)
}

func (userAPI *UserAPI) GetUserById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		Abort(context, http.StatusBadRequest, "id is required")
		return
	}
	user, err := userAPI.userRepository.GetUserById(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			Abort(context, http.StatusNotFound, "user not found")
			return
		}
		AbortWithContextError(context, http.StatusInternalServerError, "error getting user", err)
		return
	}
	context.JSON(http.StatusOK, user)
}

func (userAPI *UserAPI) CreateUser(context *gin.Context) {
	user := new(model.PostUser)
	err := context.ShouldBindJSON(user)
	if err != nil {
		Abort(context, http.StatusBadRequest, "invalid request")
		return
	}
	created, err := userAPI.userRepository.Save(user)
	if err != nil {
		AbortWithContextError(context, http.StatusInternalServerError, "error saving user", err)
		return
	}
	context.JSON(http.StatusCreated, created)
}

func (userAPI *UserAPI) DeleteUser(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		Abort(context, http.StatusBadRequest, "id is required")
		return
	}
	err := userAPI.userRepository.Delete(id)
	if err != nil {
		AbortWithContextError(context, http.StatusInternalServerError, "error deleting user", err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (userAPI *UserAPI) UpdateUser(context *gin.Context) {
	id := context.Param("id")
	exists, err := userAPI.userRepository.Exists(id)
	if err != nil {
		AbortWithContextError(context, http.StatusInternalServerError, "error updating user", err)
		return
	}
	if !exists {
		Abort(context, http.StatusNotFound, "user not found")
		return
	}
	user := new(model.PostUser)
	err = context.ShouldBindJSON(user)
	if err != nil {
		Abort(context, http.StatusBadRequest, "invalid request")
		return
	}
	updated, err := userAPI.userRepository.Update(id, user)
	if err != nil {
		AbortWithContextError(context, http.StatusInternalServerError, "error updating user", err)
		return
	}
	context.JSON(http.StatusOK, updated)
}
