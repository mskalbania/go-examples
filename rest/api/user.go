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
	Save(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
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
		abortWithError(context, http.StatusInternalServerError, "error getting users", err)
		return
	}
	context.JSON(http.StatusOK, users)
}

func (userAPI *UserAPI) GetUserById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		abortWithError(context, http.StatusBadRequest, "id is required", nil)
	}
	user, err := userAPI.userRepository.GetUserById(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			abortWithError(context, http.StatusNotFound, "user not found", err)
			return
		}
		abortWithError(context, http.StatusInternalServerError, "error getting user", err)
		return
	}
	context.JSON(http.StatusOK, user)
}

func (userAPI *UserAPI) CreateUser(context *gin.Context) {
	var user model.User
	err := context.Bind(&user)
	if err != nil {
		abortWithError(context, http.StatusBadRequest, "error reading request", err)
		return
	}
	_, err = userAPI.userRepository.Save(&user)
	if err != nil {
		abortWithError(context, http.StatusInternalServerError, "error saving user", err)
		return
	}
	context.JSON(http.StatusCreated, user)
}

func (userAPI *UserAPI) DeleteUser(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		abortWithError(context, http.StatusBadRequest, "id is required", nil)
	}
	err := userAPI.userRepository.Delete(id)
	if err != nil {
		abortWithError(context, http.StatusInternalServerError, "error deleting user", err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (userAPI *UserAPI) UpdateUser(context *gin.Context) {
	var user model.User
	err := context.ShouldBindUri(&user) //automatically binds path "id" to struct field "ID"
	if err != nil {
		abortWithError(context, http.StatusBadRequest, "uuid id is required", err)
		return
	}
	err = context.ShouldBindJSON(&user) //and now other parts of the request body
	if err != nil {
		abortWithError(context, http.StatusBadRequest, "error reading request", err)
		return
	}
	exists, err := userAPI.userRepository.Exists(user.ID)
	if err != nil {
		abortWithError(context, http.StatusInternalServerError, "error updating user", err)
		return
	}
	if exists {
		_, err := userAPI.userRepository.Update(&user)
		if err != nil {
			abortWithError(context, http.StatusInternalServerError, "error updating user", err)
			return
		}
		context.JSON(http.StatusOK, user)
		return
	}
	abortWithError(context, http.StatusNotFound, "user not found", nil)
}
