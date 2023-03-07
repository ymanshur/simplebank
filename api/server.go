package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/ymanshur/simplebank/db/sqlc"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// TODO: add routes to router

	server.router = router
	return server
}
