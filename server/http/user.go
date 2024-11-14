package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	sm "github.com/maliByatzes/socialmedia"
)

func (s *Server) addUser() gin.HandlerFunc {
	var req struct {
		User struct {
			Name           string `json:"name" binding:"required,min=3"`
			Email          string `json:"email" binding:"required,email"`
			Password       string `json:"password" binding:"required,min=8,max=72"`
			IsConsentGiven string `json:"is_consent_given" binding:"required"`
		} `json:"user" binding:"required"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		isConsentGiven, err := strconv.ParseBool(req.User.IsConsentGiven)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		newUser := sm.User{
			Name:  req.User.Name,
			Email: req.User.Email,
		}

		avatar := fmt.Sprintf("https://avatar.iran.liara.run/public/?name=%s", req.User.Name)

		emailDomain := strings.Split(req.User.Email, "@")[0]
		var role string
		if emailDomain == "mod.socialmedia.com" {
			role = "moderator"
		} else {
			role = "general"
		}

		if err := newUser.SetPassword(req.User.Password); err != nil {
			log.Printf("ERROR <addUser> - setting password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		newUser.Role = role
		newUser.Avatar = avatar

		if err := s.UserService.CreateUser(c.Request.Context(), &newUser); err != nil {
			if sm.ErrorCode(err) == sm.ECONFLICT {
				c.JSON(http.StatusConflict, gin.H{
					"error": sm.ErrorMessage(err),
				})
				return
			}

			log.Printf("ERROR <addUser> - creating new user on db: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		if !isConsentGiven {
			c.JSON(http.StatusCreated, gin.H{
				"message": "User added successfully w/o consent",
				"user":    newUser,
			})
			return
		}

		c.Set("email", req.User.Email)
		c.Set("name", req.User.Name)

		c.JSON(http.StatusCreated, gin.H{
			"message": "User added successfully",
			"user":    newUser,
		})
	}
}

func (s *Server) signin() gin.HandlerFunc {
	var req struct {
		User struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		} `json:"user" binding:"required"`
	}

	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := s.UserService.Authenticate(c.Request.Context(), req.User.Email, req.User.Password)
		if err != nil || user == nil {
			if sm.ErrorCode(err) == sm.ENOTAUTHORIZED {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": sm.ErrorMessage(err),
				})
				return
			}

			log.Printf("ERROR <signin> - authenticating the user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// NOTE: Implement auth context

		accessToken, _, err := s.TokenMaker.CreateToken(
			user.ID,
			user.Name,
			time.Hour*6,
		)
		if err != nil {
			log.Printf("ERROR <signin> - creating access token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		refreshToken, _, err := s.TokenMaker.CreateToken(
			user.ID,
			user.Name,
			time.Hour*168,
		)
		if err != nil {
			log.Printf("ERROR <signin> - creating refresh token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		if err := s.TokenService.CreateToken(c.Request.Context(), &sm.Token{
			UserID:       user.ID,
			RefreshToken: refreshToken,
			AccessToken:  accessToken,
		}); err != nil {
			log.Printf("ERROR <signin> - creating new token on db: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.Set("email", req.User.Email)
		c.Set("name", user.Name)

		c.JSON(http.StatusOK, gin.H{
			"access_token":            accessToken,
			"refresh_token":           refreshToken,
			"access_token_updated_at": time.Now(),
			"user":                    user,
		})
	}
}

func (s *Server) getCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := sm.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		// TODO: Return more info about the user here...

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func (s *Server) logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string

		if v := c.GetHeader("Authorization"); strings.HasPrefix(v, "Bearer ") {
			accessToken = strings.TrimPrefix(v, "Bearer ")
		}

		if accessToken != "" {
			tk, _, err := s.TokenService.FindTokens(c.Request.Context(), sm.TokenFilter{AccessToken: &accessToken})
			if err != nil {
				log.Printf("ERROR <logout> - finding token from db: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				return
			} else if len(tk) == 0 {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Token not found",
				})
				return
			}

			if err := s.TokenService.DeleteToken(c.Request.Context(), tk[0].ID); err != nil {
				if sm.ErrorCode(err) == sm.ENOTAUTHORIZED {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": sm.ErrorMessage(err),
					})
					return
				}

				log.Printf("ERROR <logout> - deleting token from db: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Logout successful",
			})
		}
	}
}
