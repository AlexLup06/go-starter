package middleware

import (
	"net/http"

	"alexlupatsiy.com/personal-website/backend/helpers/cookie"
	"alexlupatsiy.com/personal-website/backend/repository"
	"alexlupatsiy.com/personal-website/backend/service"
	"github.com/gin-gonic/gin"
)

func EnsureLoggedIn(sessionService *service.SessionService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ensureLoggedIn(ctx, sessionService)
		ctx.Next()
	}
}

func ensureLoggedIn(ctx *gin.Context, sessionService *service.SessionService) {
	accessCookieString, err := ctx.Cookie(repository.ACCESS_COOKIE.Type)

	// Case: We have no access cookie, so we check refresh cookie and generate new access token/cookie
	if err != nil || accessCookieString == "" {
		// check refresh token
		refreshCookieString, err := ctx.Cookie(repository.REFRESH_COOKIE.Type)
		if err != nil || refreshCookieString == "" {
			ctx.String(http.StatusUnauthorized, "Refresh Cookie not Set")
			ctx.Abort()
			return
		}

		// check refresh token against the database
		isExpired, userId, err := sessionService.VerifyRefreshToken(ctx.Request.Context(), refreshCookieString)
		if isExpired || err != nil {
			ctx.String(http.StatusUnauthorized, "Refresh Token is not valid")
			ctx.Abort()
			return
		}

		// refresh token alive and valid -> generate new Access Token
		accessTokenString, ttl, err := sessionService.CreateAccessToken(ctx.Request.Context(), userId)
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Error generating new Access token")
			ctx.Abort()
			return
		}
		cookie.SetAccessCookie(ctx, accessTokenString, ttl)
		ctx.Next()
		return
	}

	// Case: We have access cookie and access Token. Verify the access token.
	// If there's something wrong with it, it has been tinkered with most likely so no new one
	isExpired, err := sessionService.VerifyAccessToken(ctx.Request.Context(), accessCookieString)
	if isExpired || err != nil {
		ctx.String(http.StatusUnauthorized, "Access Token not valid")
		ctx.Abort()
		return
	}

	// add to Context that user is logged in

	ctx.Next()
}
