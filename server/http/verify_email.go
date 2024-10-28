package http

import (
	"github.com/gin-gonic/gin"
	"github.com/maliByatzes/socialmedia/config"
)

func sendVerificationEmail() gin.HandlerFunc {
  cfg, err := config.NewConfig()
  if err != nil {
    panic(err)
  }

  return func(ctx *gin.Context) {

  }
}
