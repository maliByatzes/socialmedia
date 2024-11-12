package http

func (s *Server) routes() {
	s.Router.Use(CorsMiddleware())

	apiRouter := s.Router.Group("/api/v1")
	{
		apiRouter.GET("/healthchecker", healthCheck())
		apiRouter.POST("/users/signup", s.addUser(), s.sendVerificationEmail())
		apiRouter.POST("/users/signin", s.signin(), s.sendLoginVerificationEmail())
	}
}
