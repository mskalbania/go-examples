package api

import (
	"github.com/gin-gonic/gin"
	"go-examples/rest/database"
)

type HealthAPI struct {
	database database.Database
}

func NewHealthAPI(database database.Database) *HealthAPI {
	return &HealthAPI{database: database}
}

func (healthAPI *HealthAPI) Health(ctx *gin.Context) {
	err := healthAPI.database.Ping(ctx)
	if err != nil {
		AbortWithContextError(ctx, 500, "db not reachable", err)
		return
	}
	ctx.Status(200)
}
