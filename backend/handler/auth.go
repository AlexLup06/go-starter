package handler

import (
	"fmt"

	"alexlupatsiy.com/personal-website/backend/helpers/ctxHelpers"
	"alexlupatsiy.com/personal-website/backend/helpers/render"
	"alexlupatsiy.com/personal-website/backend/service"
	"alexlupatsiy.com/personal-website/frontend/src/views/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	router      *gin.Engine
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(router *gin.Engine, authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		router:      router,
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Routes() {
	h.router.GET("/auth/login", h.loginGET)
	h.router.POST("/auth/login/:method", h.loginPOST)
	h.router.GET("/auth/signup", h.signupGET)
	h.router.POST("/auth/signup/:method", h.signupPOST)
}

func (h *AuthHandler) loginGET(ctx *gin.Context) {
	render.Render(ctx, 200, auth.Login())
}

func (h *AuthHandler) loginPOST(ctx *gin.Context) {
	method := ctx.Param("method")
	switch method {
	case "email":
		var loginUserWithEmailRequest service.LoginWithEmailRequest
		err := ctx.ShouldBind(&loginUserWithEmailRequest)
		if err != nil {
			email := ctx.PostForm("email")
			password := ctx.PostForm("password")

			ctxHelpers.SetKV(ctx, "email", email)
			ctxHelpers.SetKV(ctx, "password", password)
			ctxHelpers.SetKV(ctx, "isWrongEmail", "true")

			render.Render(ctx, 422, auth.LoginForm())
			return
		}

		err = h.authService.LoginWithEmail(ctx.Request.Context(), loginUserWithEmailRequest)
		if err != nil {
			render.Render(ctx, 422, auth.LoginForm())
			return
		}

		ctx.Status(200)
		return
	case "apple":
		ctx.Status(200)
		return
	case "google":
		ctx.Status(200)
		return
	default:
		ctx.Status(404)
		return
	}
}

func (h *AuthHandler) signupGET(c *gin.Context) {
	render.Render(c, 200, auth.Signup())
}

func (h *AuthHandler) signupPOST(ctx *gin.Context) {
	method := ctx.Param("method")
	switch method {
	case "email":
		var signUpUserRequest service.SignUpWithEmailRequest
		err := ctx.ShouldBind(&signUpUserRequest)
		if err != nil {
			password := ctx.PostForm("password")
			email := ctx.PostForm("email")

			ctxHelpers.SetKV(ctx, "email", email)
			ctxHelpers.SetKV(ctx, "password", password)
			ctxHelpers.SetKV(ctx, "isWrongEmail", "true")

			render.Render(ctx, 422, auth.SignupForm())
			return
		}

		err = h.userService.CreateUserWithEmail(ctx.Request.Context(), signUpUserRequest)
		if err != nil {
			fmt.Println(err)
			render.Render(ctx, 422, auth.SignupForm())
			return
		}
		ctx.Status(200)
	case "apple":
		ctx.Status(200)
	case "google":
		ctx.Status(200)
	default:
		ctx.Status(404)
	}
}
