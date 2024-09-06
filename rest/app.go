package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-examples/rest/api"
	"go-examples/rest/config"
	"go-examples/rest/database"
	"go-examples/rest/middleware"
	"go-examples/rest/repository"
	"log"
	"net/http"
	//_ "net/http/pprof" register pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartRestAPIExample() {
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatalf("env is required")
	}
	appConfig := config.Read(env)

	postgres, closable, err := database.NewPostgresDatabase(appConfig)
	defer closable()

	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	middleware.RegisterMetrics()
	authentication := middleware.NewAuthentication()
	userAPI := api.NewUserAPI(repository.NewUserRepository(postgres, &appConfig.DB))
	healthAPI := api.NewHealthAPI(postgres)

	router := setupRouter(authentication, healthAPI, userAPI)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port),
		Handler: router.Handler(),
	}

	go func() {
		//http.ErrServerClosed is returned when server.Shutdown is called
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting server: %v", err)
		}
	}()
	gracefulShutdown(srv)
}

func gracefulShutdown(server *http.Server) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown //blocks until shutdown signal is received
	log.Println("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}
	select {
	case <-ctx.Done():
		log.Println("server shutdown timeout")
	default:
	}
}

func setupRouter(auth middleware.Authentication, health api.HealthAPI, user api.UserAPI) *gin.Engine {
	g := gin.Default()
	g.Use(middleware.Metrics())

	/*Example how to wire in http profiler into gin
	g.GET("/debug/pprof/profile", gin.WrapH(http.DefaultServeMux))
	g.GET("/debug/pprof/symbol", gin.WrapH(http.DefaultServeMux))
	g.GET("/debug/pprof/trace", gin.WrapH(http.DefaultServeMux))
	g.GET("/debug/pprof/cmdline", gin.WrapH(http.DefaultServeMux))
	g.GET("/debug/pprof/", gin.WrapH(http.DefaultServeMux))
	*/

	g.GET("/metrics", middleware.MetricsHandler())
	g.GET("/health", health.Health)

	userGroup := g.Group("/api/v1").Use(auth.RequireAPIToken())
	{
		userGroup.GET("/users", user.GetUsers)
		userGroup.GET("/users/:id", user.GetUserById)
		userGroup.POST("/users", user.CreateUser)
		userGroup.DELETE("/users/:id", user.DeleteUser)
		userGroup.PUT("/users/:id", user.UpdateUser)
	}
	return g
}
