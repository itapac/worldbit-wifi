package routing

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/strusty/worldbit-wifi/services/auth"
	"github.com/strusty/worldbit-wifi/services/captcha"
	"github.com/strusty/worldbit-wifi/services/twilio"
)

type AuthRouter struct {
	authService    auth.Auth
	captchaService captcha.Captcha
	twilioService  twilio.Twilio
}

func NewAuthRouter(
	authService auth.Auth,
	captchaService captcha.Captcha,
	twilioService twilio.Twilio,
) AuthRouter {
	return AuthRouter{
		authService:    authService,
		captchaService: captchaService,
		twilioService:  twilioService,
	}
}

func (router AuthRouter) Register(group *echo.Group) {
	group.POST("/sendCode", router.sendCode)
	group.POST("", router.authenticate)
}

func (router AuthRouter) sendCode(context echo.Context) error {
	request := new(SendCodeRequest)
	if err := context.Bind(request); err != nil {
		return err
	}

	code, err := router.authService.CreateCode(request.PhoneNumber)
	if err != nil {
		return err
	}

	if err := router.twilioService.SendConfirmationCode(request.PhoneNumber, code); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}

func (router AuthRouter) authenticate(context echo.Context) error {
	request := new(VerifyCodeRequest)
	if err := context.Bind(request); err != nil {
		return err
	}

	if err := router.authService.VerifyCode(request.ConfirmationCode); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}
