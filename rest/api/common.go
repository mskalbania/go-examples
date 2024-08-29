package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-examples/rest/model"
)

func AbortWithError(context *gin.Context, status int, message string, err error) {
	context.JSON(status, model.NewError(message))
	if err == nil {
		context.Error(fmt.Errorf("%s", message))
	} else {
		context.Error(fmt.Errorf("%s: %w", message, err))
	}
	context.Abort()
}

func Abort(context *gin.Context, status int, message string) {
	context.JSON(status, model.NewError(message))
	context.Abort()
}
