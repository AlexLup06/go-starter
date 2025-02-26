package handler

import (
	"log"
	"net/http"

	"alexlupatsiy.com/personal-website/backend/helpers/cookie"
	"alexlupatsiy.com/personal-website/backend/helpers/ctxHelpers"
	customErrors "alexlupatsiy.com/personal-website/backend/helpers/errors"
	"alexlupatsiy.com/personal-website/backend/helpers/render"
	"alexlupatsiy.com/personal-website/backend/repository"
	"alexlupatsiy.com/personal-website/backend/service"
	"alexlupatsiy.com/personal-website/frontend/src/views/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	router               *gin.Engine
	authService          *service.AuthService
	userService          *service.UserService
	sessionService       *service.SessionService
	passwordResetService *service.PasswordResetService
}

func NewAuthHandler(
	router *gin.Engine,
	authService *service.AuthService,
	userService *service.UserService,
	sessionService *service.SessionService,
	passwordResetService *service.PasswordResetService,
) *AuthHandler {
	return &AuthHandler{
		router:               router,
		authService:          authService,
		userService:          userService,
		sessionService:       sessionService,
		passwordResetService: passwordResetService,
	}
}

func (h *AuthHandler) Routes(dbHandleMiddleware gin.HandlerFunc) {
	authRouter := h.router.Group("/auth", dbHandleMiddleware)

	authRouter.GET("/login", h.loginGET)
	authRouter.POST("/login/:method", h.loginPOST)
	authRouter.GET("/signup", h.signupGET)
	authRouter.POST("/signup/:method", h.signupPOST)
	authRouter.POST("/logout", h.logoutPOST)
	authRouter.GET("/request-password-reset", h.requestPasswordResetGET)
	authRouter.POST("/request-password-reset", h.requestPasswordResetPOST)
	authRouter.GET("/reset-password", h.resetPasswordGET)
	authRouter.POST("/reset-password", h.resetPasswordPOST)
	authRouter.GET("/successfull-password-reset", h.successfullPasswordResetGET)
}

func (h *AuthHandler) loginGET(ctx *gin.Context) {
	weekLoggedIn := ctxHelpers.IsWeekLoggedInCtx(ctx.Request.Context())
	if weekLoggedIn {
		ctx.Redirect(http.StatusFound, "/private")
		return
	}
	render.Render(ctx, 200, auth.Login())
}

func (h *AuthHandler) loginPOST(ctx *gin.Context) {
	method := ctx.Param("method")
	switch method {
	case "email":
		h.loginWithEmail(ctx)
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

func (h *AuthHandler) loginWithEmail(ctx *gin.Context) {
	var loginUserWithEmailRequest service.LoginWithEmailRequest
	err := ctx.ShouldBind(&loginUserWithEmailRequest)
	if err != nil {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		ctxHelpers.SetKV(ctx, "email", email)
		ctxHelpers.SetKV(ctx, "password", password)
		ctxHelpers.SetKV(ctx, "isWrongEmail", "true")

		render.Render(ctx, http.StatusUnprocessableEntity, auth.LoginForm())
		return
	}

	// login user
	userInfo, err := h.authService.LoginWithEmail(ctx.Request.Context(), loginUserWithEmailRequest)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.LoginForm())
		return
	}

	// create Refresh Token
	refreshToken, ttl, err := h.sessionService.CreateRefreshToken(ctx.Request.Context(), userInfo)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.LoginForm())
		return
	}
	// ctx.SetCookie(repository.ACCESS_COOKIE.Type, refreshToken, int(ttl), "", "", true, true)
	cookie.SetCookie(ctx, refreshToken, repository.REFRESH_COOKIE, ttl)

	// create Access Token
	accessToken, ttl, err := h.sessionService.CreateAccessToken(ctx.Request.Context(), userInfo)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.LoginForm())
		return
	}
	// ctx.SetCookie(repository.ACCESS_COOKIE.Type, accessToken, int(ttl), "", "", true, true)
	cookie.SetCookie(ctx, accessToken, repository.ACCESS_COOKIE, ttl)

	ctx.Writer.Header().Set("HX-Redirect", "/private")
	ctx.Status(http.StatusNoContent)
	return
}

func (h *AuthHandler) signupGET(ctx *gin.Context) {
	weekLoggedIn := ctxHelpers.IsWeekLoggedInCtx(ctx.Request.Context())
	if weekLoggedIn {
		ctx.Writer.Header().Set("HX-Redirect", "/private")
		ctx.Status(http.StatusNoContent)
		return
	}
	render.Render(ctx, 200, auth.Signup())
}

func (h *AuthHandler) signupPOST(ctx *gin.Context) {
	method := ctx.Param("method")
	switch method {
	case "email":
		h.signupWithEmail(ctx)
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

func (h *AuthHandler) signupWithEmail(ctx *gin.Context) {
	var signUpUserRequest service.SignUpWithEmailRequest
	err := ctx.ShouldBind(&signUpUserRequest)
	if err != nil {
		password := ctx.PostForm("password")
		email := ctx.PostForm("email")

		ctxHelpers.SetKV(ctx, "email", email)
		ctxHelpers.SetKV(ctx, "password", password)
		ctxHelpers.SetKV(ctx, "isWrongEmail", "true")

		render.Render(ctx, http.StatusUnprocessableEntity, auth.SignupForm())
		return
	}

	userId, err := h.userService.CreateUserWithEmail(ctx.Request.Context(), signUpUserRequest)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.SignupForm())
		return
	}

	// create Refresh Token
	refreshToken, ttl, err := h.sessionService.CreateRefreshToken(ctx.Request.Context(), userId)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.SignupForm())
		return
	}
	cookie.SetCookie(ctx, refreshToken, repository.REFRESH_COOKIE, ttl)

	// create Access Token
	accessToken, ttl, err := h.sessionService.CreateAccessToken(ctx.Request.Context(), userId)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.LoginForm())
		return
	}
	cookie.SetCookie(ctx, accessToken, repository.ACCESS_COOKIE, ttl)

	ctx.Writer.Header().Set("HX-Redirect", "/private")
	ctx.Status(200)
}

func (h *AuthHandler) logoutPOST(ctx *gin.Context) {
	// We have no access cookie, so we check refresh cookie and generate new access token/cookie
	// check refresh token
	refreshCookieString, err := ctx.Cookie(repository.REFRESH_COOKIE.Type)
	if err != nil || refreshCookieString == "" {
		ctx.String(http.StatusUnauthorized, "Refresh Cookie not Set")
		ctx.Abort()
		return
	}

	// check refresh token against the database
	isValid, user, err := h.sessionService.VerifyRefreshToken(ctx.Request.Context(), refreshCookieString)
	if !isValid || err != nil {
		ctx.String(http.StatusUnauthorized, "Refresh Token is not valid")
		ctx.Abort()
		return
	}

	err = h.sessionService.RevokeAllSessions(ctx.Request.Context(), user.UserId)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
	cookie.DeleteCookie(ctx, repository.ACCESS_COOKIE)
	cookie.DeleteCookie(ctx, repository.REFRESH_COOKIE)
	ctx.Writer.Header().Set("HX-Redirect", "/")
}

func (h *AuthHandler) requestPasswordResetGET(ctx *gin.Context) {
	render.Render(ctx, 200, auth.RequestPasswordReset())
}

func (h *AuthHandler) requestPasswordResetPOST(ctx *gin.Context) {
	var requestPasswordResetRequest service.RequestPasswordResetRequest
	err := ctx.ShouldBind(&requestPasswordResetRequest)
	if err != nil {
		// Wrong request
		render.Render(ctx, http.StatusUnprocessableEntity, auth.RequestPasswordReset())
		return
	}

	resetToken, err := h.passwordResetService.GenerateResetToken(ctx.Request.Context(), requestPasswordResetRequest.Email)
	if err != nil {

		if err == customErrors.ErrGeneratedTooManyResetTokens {
			log.Println("Too many password reset requests")
		}

		render.Render(ctx, 200, auth.LinkSentConfirmation())
		return
	}

	// send email with token; if email does not exist we don't care. We just don't send an email.
	// no unecessary information for bad actors that want to try out emails
	h.passwordResetService.SendResetPasswordEmailAsync(ctx.Request.Context(), resetToken)

	render.Render(ctx, 200, auth.LinkSentConfirmation())
}

func (h *AuthHandler) resetPasswordGET(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.Redirect(http.StatusFound, "/404")
		return
	}
	render.Render(ctx, 200, auth.ResetPassword(token))
}

func (h *AuthHandler) resetPasswordPOST(ctx *gin.Context) {
	var resetPasswordRequest service.ResetPasswordRequest
	err := ctx.ShouldBind(&resetPasswordRequest)
	if err != nil {
		return
	}

	// check token
	token, err := h.passwordResetService.CheckResetPasswordToken(ctx.Request.Context(), resetPasswordRequest)
	if err != nil {
		// TODO: Error handling
		if err == customErrors.ErrPasswordResetTokenNotValid {
			// if not correct then tell them to request another request-password-reset
		}
		render.Render(ctx, 200, auth.ResetPassword(resetPasswordRequest.Token))
		return
	}

	// mark token as used (just revoke all even though only this one should be marked as NOT used)
	err = h.passwordResetService.RevokeAllTokens(ctx.Request.Context(), token.Extra.UserId)

	// update password
	err = h.authService.UpdateUserPassword(ctx.Request.Context(), token.Extra.UserId, resetPasswordRequest)
	if err != nil {
		render.Render(ctx, 200, auth.ResetPassword(resetPasswordRequest.Token))
		return
	}

	// successful password reset
	ctx.Writer.Header().Set("HX-Redirect", "/auth/successfull-password-reset")
	ctx.Status(http.StatusNoContent)
	return
}

func (h *AuthHandler) successfullPasswordResetGET(ctx *gin.Context) {
	render.Render(ctx, 200, auth.SuccessfullPasswordReset())
}
