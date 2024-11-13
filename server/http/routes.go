package http

import "github.com/gin-gonic/gin"

func (s *Server) routes() {
	s.Router.Use(CorsMiddleware())

	apiRouter := s.Router.Group("/api/v1")
	{
		apiRouter.GET("/healthchecker", healthCheck())

		apiRouter.POST("/users/signup", gin.HandlerFunc(func(c *gin.Context) {
			s.addUser()(c)
			s.sendVerificationEmail()(c)
		}))
		apiRouter.POST("/users/signin", gin.HandlerFunc(func(c *gin.Context) {
			s.signin()(c)
			s.sendLoginVerificationEmail()(c)
		}))
	}
}
