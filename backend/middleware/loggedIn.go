package middleware

import (
	"net/http"

	"alexlupatsiy.com/personal-website/backend/helpers/cookie"
	"alexlupatsiy.com/personal-website/backend/helpers/ctxHelpers"
	"alexlupatsiy.com/personal-website/backend/repository"
	"alexlupatsiy.com/personal-website/backend/service"
	"github.com/gin-gonic/gin"
)

func EnsureLoggedIn(sessionService *service.SessionService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessCookieString, err := ctx.Cookie(repository.ACCESS_COOKIE.Type)

		// Case: We have no access cookie, so we check refresh cookie and generate new access token/cookie
		if err != nil || accessCookieString == "" {
			// check refresh token
			refreshCookieString, err := ctx.Cookie(repository.REFRESH_COOKIE.Type)
			if err != nil || refreshCookieString == "" {
				ctx.Redirect(http.StatusFound, "/auth/login")
				ctx.Abort()
				return
			}

			// check refresh token against the database
			isValid, user, err := sessionService.VerifyRefreshToken(ctx.Request.Context(), refreshCookieString)
			if !isValid || err != nil {
				ctx.Redirect(http.StatusFound, "/auth/login")
				ctx.Abort()
				return
			}

			// refresh token alive and valid -> generate new Access Token
			accessTokenString, ttl, err := sessionService.CreateAccessToken(ctx.Request.Context(), user)
			if err != nil {
				ctx.Redirect(http.StatusFound, "/auth/login")
				ctx.Abort()
				return
			}
			cookie.SetCookie(ctx, accessTokenString, repository.ACCESS_COOKIE, ttl)
			ctx.Next()
			return
		}

		// Case: We have access cookie and access Token. Verify the access token.
		// If there's something wrong with it, it has been tinkered with most likely so no new one
		isValid, userId, err := sessionService.VerifyAccessToken(ctx.Request.Context(), accessCookieString)
		if !isValid || err != nil {
			// ctx.String(http.StatusUnauthorized, "Access Token not valid")
			ctx.Redirect(http.StatusFound, "/auth/login")
			ctx.Abort()
			return
		}

		loggedInCtx := ctxHelpers.WithIsLoggedInCtx(ctx.Request.Context())
		if userId != "" {
			loggedInCtx = ctxHelpers.WithUserIdCtx(loggedInCtx, userId)
		}
		ctx.Request = ctx.Request.WithContext(loggedInCtx)

		ctx.Next()

	}
}

// func ensureLoggedIn(ctx *gin.Context, sessionService *service.SessionService) {
// 	accessCookieString, err := ctx.Cookie(repository.ACCESS_COOKIE.Type)

// 	// Case: We have no access cookie, so we check refresh cookie and generate new access token/cookie
// 	if err != nil || accessCookieString == "" {
// 		// check refresh token
// 		refreshCookieString, err := ctx.Cookie(repository.REFRESH_COOKIE.Type)
// 		if err != nil || refreshCookieString == "" {
// 			ctx.Redirect(http.StatusFound, "/auth/login")
// 			ctx.Abort()
// 			return
// 		}

// 		// check refresh token against the database
// 		isValid, user, err := sessionService.VerifyRefreshToken(ctx.Request.Context(), refreshCookieString)
// 		if !isValid || err != nil {
// 			// ctx.String(http.StatusUnauthorized, "Refresh Token is not valid")
// 			ctx.Redirect(http.StatusFound, "/auth/login")
// 			ctx.Abort()
// 			return
// 		}

// 		// refresh token alive and valid -> generate new Access Token
// 		accessTokenString, ttl, err := sessionService.CreateAccessToken(ctx.Request.Context(), user)
// 		if err != nil {
// 			// ctx.String(http.StatusUnauthorized, "Error generating new Access token")
// 			ctx.Redirect(http.StatusFound, "/auth/login")
// 			ctx.Abort()
// 			return
// 		}
// 		cookie.SetCookie(ctx, accessTokenString, repository.ACCESS_COOKIE, ttl)
// 		ctx.Next()
// 		return
// 	}

// 	// Case: We have access cookie and access Token. Verify the access token.
// 	// If there's something wrong with it, it has been tinkered with most likely so no new one
// 	isValid, userId, err := sessionService.VerifyAccessToken(ctx.Request.Context(), accessCookieString)
// 	if !isValid || err != nil {
// 		// ctx.String(http.StatusUnauthorized, "Access Token not valid")
// 		ctx.Redirect(http.StatusFound, "/auth/login")
// 		ctx.Abort()
// 		return
// 	}

// 	loggedInCtx := ctxHelpers.WithIsLoggedInCtx(ctx.Request.Context())
// 	if userId != "" {
// 		loggedInCtx = ctxHelpers.WithUserIdCtx(loggedInCtx, userId)
// 	}
// 	ctx.Request = ctx.Request.WithContext(loggedInCtx)

// 	ctx.Next()
// }
