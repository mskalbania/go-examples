package auth

import (
	"github.com/gin-gonic/gin"
	"go-examples/rest/api"
	"net/http"
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
			api.Abort(context, http.StatusUnauthorized, "missing api key")
			return
		}
		if !auth.allowlist[apiKey] {
			api.Abort(context, http.StatusUnauthorized, "invalid api key")
			return
		}
	}
}
