package handler

import (
	"fmt"
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
	authRouter.POST("/login", h.loginWithEmailPOST)
	authRouter.POST("/apple", h.signinWithApplePOST)
	authRouter.POST("/google", h.signinWithGooglePOST)

	authRouter.GET("/signup", h.signupGET)
	authRouter.POST("/signup", h.signupWithEmailPOST)
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
		fmt.Println("we are week logged in")
		ctx.Redirect(http.StatusFound, "/private")
		return
	}
	render.Render(ctx, 200, auth.Login())
}

func (h *AuthHandler) loginWithEmailPOST(ctx *gin.Context) {
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
	cookie.SetCookie(ctx, refreshToken, repository.REFRESH_COOKIE, ttl)

	// create Access Token
	accessToken, ttl, err := h.sessionService.CreateAccessToken(ctx.Request.Context(), userInfo)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.LoginForm())
		return
	}
	cookie.SetCookie(ctx, accessToken, repository.ACCESS_COOKIE, ttl)

	ctx.Writer.Header().Set("HX-Redirect", "/private")
	ctx.Status(http.StatusNoContent)
	return
}

func (h *AuthHandler) signinWithApplePOST(ctx *gin.Context) {
	// var authResponse service.AppleAuthResponse
	// err := ctx.ShouldBind(&authResponse)
	// // Bind Apple response to struct
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	// 	return
	// }

	// // Step 1: Verify the ID Token (JWT)
	// claims, err := h.authService.VerifyAppleIDToken(authResponse.IDToken)
	// if err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
	// 	return
	// }

	// // Step 2: Exchange Code for Apple Tokens (Access/Refresh)
	// appleTokens, err := h.authService.ExchangeAppleCodeForTokens(authResponse.Code)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange authorization code"})
	// 	return
	// }

	// // Step 3: Handle user authentication (create/find user in DB)
	// userID := claims["sub"].(string)
	// email := claims["email"].(string)

	// // Find user in DB, or create a new one if needed
	// user := h.userService.CreateUserWithApple(userID, email, appleTokens.RefreshToken)

	// // Step 4: Issue Your Own Tokens (Access & Refresh)
	// accessToken, refreshToken := GenerateAppTokens(user.ID)

	// // Step 5: Set tokens in HTTP-only cookies
	// c.SetCookie("access_token", accessToken, 900, "/", "", true, true)       // 15 mins
	// c.SetCookie("refresh_token", refreshToken, 1209600, "/", "", true, true) // 14 days

}

func (h *AuthHandler) signinWithGooglePOST(ctx *gin.Context) {
	var signInWithGoogle service.SignInWithGoogle
	err := ctx.ShouldBind(&signInWithGoogle)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.Login())
		return
	}

	// CSRF Check
	csrfCookie, err := ctx.Request.Cookie("g_csrf_token")
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.Login())
		return
	}

	if signInWithGoogle.CSRFToken == "" || signInWithGoogle.CSRFToken != csrfCookie.Value {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.Login())
		return
	}

	// 3. Validate the ID token
	payload, err := h.authService.ValidateGoogleIdToken(ctx.Request.Context(), signInWithGoogle.IdToken)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.Login())
		return
	}

	// Create new user or login
	userInfo, err := h.authService.GoogleLogin(ctx.Request.Context(), payload)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.Login())
		return
	}

	// create Refresh Token
	refreshToken, ttl, err := h.sessionService.CreateRefreshToken(ctx.Request.Context(), userInfo)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.Login())
		return
	}
	cookie.SetCookie(ctx, refreshToken, repository.REFRESH_COOKIE, ttl)

	// create Access Token
	accessToken, ttl, err := h.sessionService.CreateAccessToken(ctx.Request.Context(), userInfo)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.Login())
		return
	}
	cookie.SetCookie(ctx, accessToken, repository.ACCESS_COOKIE, ttl)

	ctx.Redirect(http.StatusFound, "/private")
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

func (h *AuthHandler) signupWithEmailPOST(ctx *gin.Context) {
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

	userInfo, err := h.authService.EmailSignUp(ctx.Request.Context(), signUpUserRequest)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.SignupForm())
		return
	}

	// create Refresh Token
	refreshToken, ttl, err := h.sessionService.CreateRefreshToken(ctx.Request.Context(), userInfo)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.SignupForm())
		return
	}
	cookie.SetCookie(ctx, refreshToken, repository.REFRESH_COOKIE, ttl)

	// create Access Token
	accessToken, ttl, err := h.sessionService.CreateAccessToken(ctx.Request.Context(), userInfo)
	if err != nil {
		render.Render(ctx, http.StatusUnprocessableEntity, auth.SignupForm())
		return
	}
	cookie.SetCookie(ctx, accessToken, repository.ACCESS_COOKIE, ttl)

	ctx.Writer.Header().Set("HX-Redirect", "/private")
	ctx.Status(200)
}

func (h *AuthHandler) logoutPOST(ctx *gin.Context) {
	// We have no access cookie, so we check refresh cookie and generate new access token/cookie
	// check refresh token
	cookie.DeleteCookie(ctx, repository.ACCESS_COOKIE)
	cookie.DeleteCookie(ctx, repository.REFRESH_COOKIE)
	ctx.Writer.Header().Set("HX-Redirect", "/")

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
}

func (h *AuthHandler) requestPasswordResetGET(ctx *gin.Context) {
	weekLoggedIn := ctxHelpers.IsWeekLoggedInCtx(ctx.Request.Context())
	if weekLoggedIn {
		ctx.Writer.Header().Set("HX-Redirect", "/")
		ctx.Status(http.StatusNoContent)
		return
	}
	render.Render(ctx, 200, auth.RequestPasswordReset())
}

func (h *AuthHandler) requestPasswordResetPOST(ctx *gin.Context) {
	weekLoggedIn := ctxHelpers.IsWeekLoggedInCtx(ctx.Request.Context())
	if weekLoggedIn {
		ctx.Writer.Header().Set("HX-Redirect", "/")
		ctx.Status(http.StatusNoContent)
		return
	}
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
	weekLoggedIn := ctxHelpers.IsWeekLoggedInCtx(ctx.Request.Context())
	if weekLoggedIn {
		ctx.Writer.Header().Set("HX-Redirect", "/")
		ctx.Status(http.StatusNoContent)
		return
	}
	token := ctx.Query("token")
	if token == "" {
		ctx.Redirect(http.StatusFound, "/404")
		return
	}
	render.Render(ctx, 200, auth.ResetPassword(token))
}

func (h *AuthHandler) resetPasswordPOST(ctx *gin.Context) {
	weekLoggedIn := ctxHelpers.IsWeekLoggedInCtx(ctx.Request.Context())
	if weekLoggedIn {
		ctx.Writer.Header().Set("HX-Redirect", "/")
		ctx.Status(http.StatusNoContent)
		return
	}
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
