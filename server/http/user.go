package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

		newUser.SetPassword(req.User.Password)
		newUser.Role = role
		newUser.Avatar = avatar

		if err := s.UserService.CreateUser(c.Request.Context(), &newUser); err != nil {
			if sm.ErrorCode(err) == sm.ECONFLICT {
				c.JSON(http.StatusConflict, gin.H{
					"error": sm.ErrorMessage(err),
				})
				return
			}

			log.Printf("error in create user handler: %v", err)
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
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"message": "User added successfully w/ consent",
				"user":    newUser,
			})
			// c.Next()
		}
	}
}
