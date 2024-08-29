package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-examples/rest/database"
)

type HealthAPI struct {
	postgres *database.PostgresDatabase
}

func NewHealthAPI(postgres *database.PostgresDatabase) *HealthAPI {
	return &HealthAPI{postgres}
}

func (healthAPI *HealthAPI) Health(ctx *gin.Context) {
	err := healthAPI.postgres.Conn.Ping(context.TODO())
	if err != nil {
		abortWithError(ctx, 500, "db not reachable", err)
		return
	}
	ctx.Status(200)
}
