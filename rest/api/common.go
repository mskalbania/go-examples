package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-examples/rest/model"
)

func AbortWithContextError(context *gin.Context, status int, message string, err error) {
	context.JSON(status, model.NewError(message))
	context.Error(fmt.Errorf("%s: %w", message, err))
	context.Abort()
}

func Abort(context *gin.Context, status int, message string) {
	context.JSON(status, model.NewError(message))
	context.Abort()
}
