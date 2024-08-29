package auth

import (
	"github.com/gin-gonic/gin"
	"go-examples/rest/api"
)

var apiKeyHeader = "X-API-KEY"

type Authentication struct {
	allowlist map[string]bool
}

//Naive implementation just to show example of auth middleware

func NewAuthentication() *Authentication {
	return &Authentication{allowlist: map[string]bool{
		"token": true,
	}}
}

func (auth *Authentication) RequireAPIToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		apiKey := context.GetHeader(apiKeyHeader)
		if apiKey == "" {
			api.AbortWithError(context, 401, "missing api key", nil)
			return
		}
		if !auth.allowlist[apiKey] {
			api.AbortWithError(context, 401, "invalid api key", nil)
			return
		}
	}
}
