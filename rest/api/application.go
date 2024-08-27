package api

import (
	"errors"
	"fmt"
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
	Exists(id string) bool
	Delete(id string) error
}

type User struct {
	userRepository UserRepository
}

func NewUserAPI(userRepository UserRepository) *User {
	return &User{userRepository: userRepository}
}

func (app *User) GetUsers(context *gin.Context) {
	users, err := app.userRepository.GetAllUsers()
	if err != nil {
		abortWithError(context, http.StatusInternalServerError, "error getting users", err)
		return
	}
	context.JSON(http.StatusOK, users)
}

func (app *User) GetUserById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		abortWithError(context, http.StatusBadRequest, "id is required", nil)
	}
	user, err := app.userRepository.GetUserById(id)
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

func (app *User) CreateUser(context *gin.Context) {
	var user model.User
	err := context.Bind(&user)
	if err != nil {
		abortWithError(context, http.StatusBadRequest, "error reading request", err)
		return
	}
	_, err = app.userRepository.Save(&user)
	if err != nil {
		abortWithError(context, http.StatusInternalServerError, "error saving user", err)
		return
	}
	context.JSON(http.StatusCreated, user)
}

func (app *User) DeleteUser(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		abortWithError(context, http.StatusBadRequest, "id is required", nil)
	}
	err := app.userRepository.Delete(id)
	if err != nil {
		abortWithError(context, http.StatusInternalServerError, "error deleting user", err)
		return
	}
	context.Status(http.StatusNoContent)
}

func (app *User) UpdateUser(context *gin.Context) {
	var user model.User
	err := context.ShouldBindUri(&user) //automatically binds path "id" to struct field "ID"
	if err != nil {
		abortWithError(context, http.StatusBadRequest, "id is required", err)
		return
	}
	err = context.ShouldBindJSON(&user) //and now other parts of the request body
	if err != nil {
		abortWithError(context, http.StatusBadRequest, "error reading request", err)
		return
	}
	if app.userRepository.Exists(user.ID) {
		_, err := app.userRepository.Update(&user)
		if err != nil {
			abortWithError(context, http.StatusInternalServerError, "error updating user", err)
			return
		}
		context.JSON(http.StatusOK, user)
		return
	}
	abortWithError(context, http.StatusNotFound, "user not found", nil)
}

func abortWithError(context *gin.Context, status int, message string, err error) {
	context.JSON(status, model.NewError(message))
	context.Error(fmt.Errorf("%s: %w", message, err))
	context.Abort()
}
