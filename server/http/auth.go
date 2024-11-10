package http

import (
	"reflect"

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
