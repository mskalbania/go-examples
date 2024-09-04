package auth

import (
	"github.com/gin-gonic/gin"
	"go-examples/rest/api"
	"net/http"
)

var apiKeyHeader = "X-API-KEY"

type Authentication interface {
	RequireAPIToken() gin.HandlerFunc
}

type authentication struct {
	allowlist map[string]bool
}

//Naive implementation just to show example of auth middleware

func NewAuthentication() Authentication {
	return &authentication{allowlist: map[string]bool{
		"token": true,
	}}
}

func (auth *authentication) RequireAPIToken() gin.HandlerFunc {
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
