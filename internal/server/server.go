package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/blazeisclone/vehicle-dms-inventory/events"
	"github.com/blazeisclone/vehicle-dms-inventory/internal/database"
)

type Server struct {
	port      int
	db        database.Service
	publisher events.Publisher
}

func NewServer(pub events.Publisher) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:      port,
		db:        database.New(),
		publisher: pub,
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
