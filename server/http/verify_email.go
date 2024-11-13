package http

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"

	"github.com/gin-gonic/gin"
	sm "github.com/maliByatzes/socialmedia"
	"github.com/maliByatzes/socialmedia/config"
	"github.com/maliByatzes/socialmedia/mail"
)

func (s *Server) sendVerificationEmail() gin.HandlerFunc {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Next()

		emailAny, exists := c.Get("email")
		if !exists {
			return
		}
		nameAny, _ := c.Get("name")

		bgCtx := context.Background()

		go func() {

			email := fmt.Sprintf("%v", emailAny)
			name := fmt.Sprintf("%v", nameAny)

			verificationCode := rand.IntN(90000) + 10000
			verificationLink := fmt.Sprintf("%s/auth/verify?code=%d&email=%s",
				cfg.ClientURL, verificationCode, email)

			sender := mail.NewGmailSender("SocialMedia", cfg.Email, cfg.EmailPassword)
			content := verifyEmailHTML(string(name), verificationLink, verificationCode)
			if err := sender.SendEmail("Verify your email address", content, []string{email}, nil, nil); err != nil {
				log.Printf("error sending an email: %v", err)
				/*contextCopy.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})*/
				return
			}

			newEmailVerification := sm.Email{
				Email:            email,
				VerificationCode: fmt.Sprintf("%d", verificationCode),
				For:              "signup",
			}

			if err := s.EmailService.CreateEmailVerification(bgCtx, &newEmailVerification); err != nil {
				log.Printf("error in create email verification: %v", err)
				/*c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})*/
				return
			}

			/*c.JSON(http.StatusOK, gin.H{
				"message": "Email verification sent successfully",
			})*/
		}()
	}
}

func (s *Server) sendLoginVerificationEmail() gin.HandlerFunc {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Next()

		emailAny, exists := c.Get("email")
		if !exists {
			return
		}
		nameAny, _ := c.Get("name")

		bgCtx := c.Copy()

		go func() {

			email := fmt.Sprintf("%v", emailAny)
			name := fmt.Sprintf("%v", nameAny)

			// NOTE: implement context data
			currentContextData := sm.Context{}

			verificationLink := fmt.Sprintf("%s/verify-login?id=%d&email=%s", cfg.ClientURL, currentContextData.ID, email)
			blockLink := fmt.Sprintf("%s/block-device?id=%d&email=%s", cfg.ClientURL, currentContextData.ID, email)

			sender := mail.NewGmailSender("SocialMedia", cfg.Email, cfg.EmailPassword)
			content := verifyLoginHTML(name, verificationLink, blockLink, currentContextData)

			if err := sender.SendEmail("Action Required: Verify Recent Login", content, []string{email}, nil, nil); err != nil {
				log.Printf("error sending an email: %v", err)
				/*c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})*/
				return
			}

			newEmailVerification := sm.Email{
				Email:            email,
				VerificationCode: strconv.Itoa(int(currentContextData.ID)),
				For:              "login",
			}

			if err := s.EmailService.CreateEmailVerification(bgCtx, &newEmailVerification); err != nil {
				log.Printf("error in create email verification: %v", err)
				/*c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})*/
				return
			}

			/*c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Access blocked due to suspicious activity. Verification email was sent to your email addres.",
			})*/
		}()
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

func verifyLoginHTML(name, verificationLink, blockLink string, currentContextData sm.Context) string {
	return fmt.Sprintf(`
	<div style="background-color: #F4F4F4; padding: 20px;">
      <div style="background-color: #fff; padding: 20px; border-radius: 10px;">
        <h1 style="color: black; font-size: 24px; margin-bottom: 20px;">New login attempt detected</h1>
        <p>Dear %s,</p>
        <p>Our system has detected that a new login was attempted from the following device and location at ${currentContextData.time}:</p>
        <ul style="list-style: none; padding-left: 0;">
          <li><strong>IP Address:</strong> %s</li>
          <li><strong>Location:</strong> %s, %s</li>
          <li><strong>Device:</strong> %s %s</li>
          <li><strong>Browser:</strong> %s</li>
          <li><strong>Operating System:</strong> %s</li>
          <li><strong>Platform:</strong> %s</li>
        </ul>
        <p>If this was you, please click the button below to verify your login:</p>
        <div style="text-align: center;">
          <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #1da1f2; color: #fff; text-decoration: none; border-radius: 5px; margin-bottom: 20px;">Verify Login</a>
        </div>
        <p>If you believe this was an unauthorized attempt, please click the button below to block this login:</p>
        <div style="text-align: center;">
          <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #E0245E; color: #fff; text-decoration: none; border-radius: 5px; margin-bottom: 20px;">Block Login</a>
        </div>
        <p>Please verify that this login was authorized. If you have any questions or concerns, please contact our customer support team.</p>
      </div>
    </div>
	`, name,
		currentContextData.IP,
		currentContextData.City,
		currentContextData.Country,
		currentContextData.Device,
		currentContextData.DeviceType,
		currentContextData.Browser,
		currentContextData.OS,
		currentContextData.Platform,
		verificationLink,
		blockLink,
	)
}
