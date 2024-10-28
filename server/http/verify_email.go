package http

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maliByatzes/socialmedia/config"
	"github.com/maliByatzes/socialmedia/mail"
)

func sendVerificationEmail() gin.HandlerFunc {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		emailAny, _ := c.Get("email")
		nameAny, _ := c.Get("name")

		email := fmt.Sprintf("%v", emailAny)
		name := fmt.Sprintf("%v", nameAny)

		verificationCode := rand.IntN(90000) + 10000
		verificationLink := fmt.Sprintf("%s/auth/verify?code=%d&email=%s",
			cfg.ClientURL, verificationCode, email)

		sender := mail.NewGmailSender("SocialMedia", cfg.Email, cfg.EmailPassword)
		content := verifyEmailHTML(string(name), verificationLink, verificationCode)
		sender.SendEmail("Verify your email address", content, []string{email}, nil, nil)

		// create a new verification

		// and then return

		c.JSON(http.StatusOK, gin.H{
			"message": "Verification email successfully sent.",
		})
	}
}

func verifyEmailHTML(name, verificationLink string, verificationCode int) string {
	return fmt.Sprintf(`
  <div style="max-width: 600px; margin: auto; background-color: #f4f4f4; padding: 20px; border-radius: 10px; box-shadow: 0 2px 4px rgb(104, 182, 255);">
    <div style="background-color: #ffffff; padding: 20px; border-radius: 10px;">
      <h1 style="font-size: 24px; margin-bottom: 20px; text-align: center; color: #6AFF5E; font-weight: bold">SocialMedia</h1>
      <p style="font-size: 18px; margin-bottom: 20px; text-align: center; color: #4b5563; font-weight: bold;">Welcome to SocialMedia, %s!</p>
      <p style="font-size: 16px; margin-bottom: 20px; text-align: center; color: #4b5563;">Please click the button below to verify your email address and activate your account:</p>
      <div style="text-align: center; margin-bottom: 20px;">
        <a href="%s" style="background-color: #3b82f6; color: #ffffff; padding: 12px 25px; border-radius: 5px; text-decoration: none; display: inline-block; font-size: 16px; font-weight: bold;">Verify Email Address</a>
      </div>
      <p style="font-size: 14px; margin-bottom: 20px; text-align: center; color: #4b5563;">Please note that the device you are using for this verification process will be set as your primary device.</p>
      <p style="font-size: 14px; margin-bottom: 20px; text-align: center; color: #6b7280;">The link will expire in 30 minutes.</p>
      <p style="font-size: 16px; margin-bottom: 15px; text-align: center; color: #3b82f6; font-weight: bold;">Your verification code is: <span style="color: #000000;">%d</span></p>
      <p style="font-size: 14px; margin-bottom: 20px; text-align: center; color: #4b5563;">If you did not create an account, please ignore this email.</p>
    </div>
  </div>`, name, verificationLink, verificationCode)
}
