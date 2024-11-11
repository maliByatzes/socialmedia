package http

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	sm "github.com/maliByatzes/socialmedia"
	"github.com/maliByatzes/socialmedia/postgres"
	"github.com/maliByatzes/socialmedia/token"
)

const Timeout = 5 * time.Second

type Server struct {
	Server                 *http.Server
	Router                 *gin.Engine
	TokenMaker             token.Maker
	UserService            sm.UserService
	EmailService           sm.EmailService
	ContextService         sm.ContextService
	PreferenceService      sm.PreferenceService
	SuspiciousLoginService sm.SuspiciousLoginService
}

func NewServer(db *postgres.DB, secretKey string) (*Server, error) {
	s := Server{
		Server: &http.Server{
			WriteTimeout: Timeout,
			ReadTimeout:  Timeout,
			IdleTimeout:  Timeout,
		},
		Router: gin.Default(),
	}

	tkMaker, err := token.NewJWTMaker(secretKey)
	if err != nil {
		return nil, err
	}
	s.TokenMaker = tkMaker

	s.routes()
	s.UserService = postgres.NewUserService(db)
	s.EmailService = postgres.NewEmailService(db)
	s.ContextService = postgres.NewContextService(db)
	s.PreferenceService = postgres.NewPreferenceService(db)
	s.SuspiciousLoginService = postgres.NewSuspiciousLoginService(db)
	s.Server.Handler = s.Router

	return &s, nil
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
