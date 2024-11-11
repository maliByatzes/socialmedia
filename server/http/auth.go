package http

import (
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	sm "github.com/maliByatzes/socialmedia"
	"github.com/maliByatzes/socialmedia/utils"
)

// Helper functions
func isTrustedDevice(cD *utils.IPContext, uD *sm.Context) bool {
	return cD.IP == uD.IP &&
		cD.Country == uD.Country &&
		cD.City == uD.City &&
		cD.Browser == uD.Browser &&
		cD.Platform == uD.Platform &&
		cD.OS == uD.OS &&
		cD.Device == uD.Device &&
		cD.DeviceType == uD.DeviceType
}

func isSuspciousContextChanged(oldContextData, newContextData utils.IPContext) bool {
	return !reflect.DeepEqual(oldContextData, newContextData)
}

func isOldDataMatched(oldSuspiciouseContextData, userContextData utils.IPContext) bool {
	return reflect.DeepEqual(oldSuspiciouseContextData, userContextData)
}

// run in go routine
func (s *Server) verifyContextData(ctx *gin.Context, existingUser *sm.User) error {
	userContextData, err := s.ContextService.FindContextByID(ctx.Request.Context(), existingUser.ID)
	if err != nil {
		return err
	}

	currentContextData, err := utils.GetCurrentContextData(ctx.ClientIP(), ctx.Request)
	if err != nil {
		return err
	}

	if isTrustedDevice(currentContextData, userContextData) {
		return nil
	}

	// get Old suspcicious login context data

	return nil
}

// GET /auth/context-data/primary
func (s *Server) getAuthContextData() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := sm.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		result, err := s.ContextService.FindContextByUserID(c.Request.Context(), user.ID)
		if err != nil {
			if sm.ErrorCode(err) == sm.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			log.Printf("error in getAuthContextData: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"context": result,
		})
	}
}

// GET /auth/context-data/trusted
func (s *Server) getTrustedAuthContextData() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := sm.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		bTrue := true
		bFalse := false
		result, n, err := s.SuspiciousLoginService.FindSLs(c.Request.Context(), sm.SLFilter{
			UserID:    &user.ID,
			IsTrusted: &bTrue,
			IsBlocked: &bFalse,
		})

		if err != nil {
			log.Printf("error in getAuthContextData: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"n":        n,
			"contexts": result,
		})
	}
}

// GET /auth/context-data/blocked
func (s *Server) getBlockedAuthContextData() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := sm.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		bTrue := true
		bFalse := false
		result, n, err := s.SuspiciousLoginService.FindSLs(c.Request.Context(), sm.SLFilter{
			UserID:    &user.ID,
			IsTrusted: &bFalse,
			IsBlocked: &bTrue,
		})

		if err != nil {
			log.Printf("error in getAuthContextData: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"n":        n,
			"contexts": result,
		})
	}
}

// GET /auth/user-preferences
func (s *Server) getUserPreferences() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := sm.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		userPreference, err := s.PreferenceService.FindPreferenceByUserID(c.Request.Context(), user.ID)
		if err != nil {
			if sm.ErrorCode(err) == sm.ENOTFOUND {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}

			log.Printf("error in getUserPreferences: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_preferences": userPreference,
		})
	}
}

// DELETE /auth/context-data/:contextId
func (s *Server) deleteContextAuthData() gin.HandlerFunc {
	return func(c *gin.Context) {

		contextIDstr := c.Param("contextId")
		contextID, err := strconv.ParseUint(contextIDstr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid context id param",
			})
			return
		}

		user := sm.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		if err := s.SuspiciousLoginService.DeleteSL(c.Request.Context(), uint(contextID)); err != nil {
			log.Printf("error in deleteContextAuthData: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Data deleted successfully",
		})
	}
}

// PATCH /auth/context-data/block/:contextId
func (s *Server) blockContextAuthData() gin.HandlerFunc {
	return func(c *gin.Context) {

		contextIDstr := c.Param("contextId")
		contextID, err := strconv.ParseUint(contextIDstr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid context id param",
			})
			return
		}

		user := sm.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		bTrue := true
		bFalse := false
		_, err = s.SuspiciousLoginService.UpdateSL(c.Request.Context(), uint(contextID), sm.SLUpdate{
			IsBlocked: &bTrue,
			IsTrusted: &bFalse,
		})

		if err != nil {
			log.Printf("error in blockContextAuthData: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Blocked successfully",
		})
	}
}

// PATCH /auth/context-data/unblock/:contextId
func (s *Server) unblockContextAuthData() gin.HandlerFunc {
	return func(c *gin.Context) {

		contextIDstr := c.Param("contextId")
		contextID, err := strconv.ParseUint(contextIDstr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid context id param",
			})
			return
		}

		user := sm.UserFromContext(c.Request.Context())
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		bTrue := true
		bFalse := false
		_, err = s.SuspiciousLoginService.UpdateSL(c.Request.Context(), uint(contextID), sm.SLUpdate{
			IsBlocked: &bFalse,
			IsTrusted: &bTrue,
		})

		if err != nil {
			log.Printf("error in blockContextAuthData: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Unblocked successfully",
		})
	}
}
