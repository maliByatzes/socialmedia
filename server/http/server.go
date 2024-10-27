package http

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maliByatzes/socialmedia/postgres"
)

const Timeout = 5 * time.Second

type Server struct {
	Server *http.Server
	Router *gin.Engine
}

func NewServer(db *postgres.DB) *Server {
	s := Server{
		Server: &http.Server{
			WriteTimeout: Timeout,
			ReadTimeout:  Timeout,
			IdleTimeout:  Timeout,
		},
		Router: gin.Default(),
	}

  s.routes()
  s.Server.Handler = s.Router

	return &s
}

func (s *Server) Run(port string) error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	s.Server.Addr = port
	log.Printf("🗿 Server is starting on port %s", port)
	return s.Server.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.Server.Shutdown(ctx)
}

func healthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	}
}
