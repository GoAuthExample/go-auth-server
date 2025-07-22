package server

import (
	"fmt"
	"goAuthExample/internal/auth"
	"goAuthExample/internal/database"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {

	auth.NewAuth()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
