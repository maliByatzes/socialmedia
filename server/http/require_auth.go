package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sm "github.com/maliByatzes/socialmedia"
)

func (s *Server) requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		var accessToken string

		if v := c.GetHeader("Authorization"); strings.HasPrefix(v, "Bearer ") {
			accessToken = strings.TrimPrefix(v, "Bearer ")
		} else {
			accessToken, _ = c.Cookie("access_token")
		}

		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - No access token",
			})
			c.Abort()
			return
		}

		payload, err := s.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("Unauthorized - %v", err),
			})
			c.Abort()
			return
		}

		user, err := s.UserService.FindUserByID(c.Request.Context(), payload.ID)
		if err != nil {
			if sm.ErrorCode(err) == sm.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
					"error": sm.ErrorMessage(err),
				})
				c.Abort()
				return
			}

			log.Printf("error in requireAuth: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			c.Abort()
			return
		}

		ctx := sm.NewContextWithUser(c.Request.Context(), user)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
