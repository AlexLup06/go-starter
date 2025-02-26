package middleware

import (
	"alexlupatsiy.com/personal-website/backend/helpers/ctxHelpers"
	"alexlupatsiy.com/personal-website/backend/repository"
	"alexlupatsiy.com/personal-website/backend/service"
	"github.com/gin-gonic/gin"
)

func SetUserInfo(sessionService *service.SessionService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshCookieString, err := ctx.Cookie(repository.REFRESH_COOKIE.Type)
		if err != nil || refreshCookieString == "" {
			ctx.Next()
			return
		}

		userInfo, err := sessionService.ParseUserInfo(refreshCookieString)
		if err != nil {
			ctx.Next()
			return
		}

		// give global access to Username and Email; for example for navbar to always show username without hard checking the refresh token
		userInfoContext := ctxHelpers.WithUsernameCtx(ctx.Request.Context(), userInfo.Username)
		userInfoContext = ctxHelpers.WithEmailCtx(userInfoContext, userInfo.Email)
		userInfoContext = ctxHelpers.WithIsWeekLoggedInCtx(userInfoContext)
		ctx.Request = ctx.Request.WithContext(userInfoContext)

		ctx.Next()
	}
}
