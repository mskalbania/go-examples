package rest

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go-examples/rest/api"
	"go-examples/rest/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//TODO:
//1. sqlLite repository
//2. Tests
//3. Add validation

func StartRestAPIExample() {
	g := gin.Default()

	app := api.NewUserAPI(repository.NewInMemoryUserRepository())

	g.Group("v1")
	{
		g.GET("/users", app.GetUsers)
		g.GET("/users/:id", app.GetUserById)
		g.POST("/users", app.CreateUser)
		g.DELETE("/users/:id", app.DeleteUser)
		g.PUT("/users/:id", app.UpdateUser)
	}

	srv := &http.Server{
		Addr:    "localhost:8080",
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

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
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
