package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-examples/rest/api"
	"go-examples/rest/config"
	"go-examples/rest/database"
	"go-examples/rest/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//TODO
//Add validation
//Tests
//Contexts

func StartRestAPIExample() {
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatalf("env is required")
	}
	appConfig := config.Read(env)

	g := gin.Default()

	postgres, closable, err := database.NewPostgresDatabase(appConfig)
	defer closable()

	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	userAPI := api.NewUserAPI(repository.NewUserRepository(postgres))
	healthAPI := api.NewHealthAPI(postgres)

	g.GET("/health", healthAPI.Health)

	g.Group("v1")
	{
		g.GET("/users", userAPI.GetUsers)
		g.GET("/users/:id", userAPI.GetUserById)
		g.POST("/users", userAPI.CreateUser)
		g.DELETE("/users/:id", userAPI.DeleteUser)
		g.PUT("/users/:id", userAPI.UpdateUser)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port),
		Handler: g.Handler(),
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
