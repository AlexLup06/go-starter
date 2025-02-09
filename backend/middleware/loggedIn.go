package middleware

import (
	"alexlupatsiy.com/personal-website/backend/service"
	"github.com/gin-gonic/gin"
)

func EnsureLoggedIn(authService *service.AuthService, sessionService *service.SessionService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ensureLoggedIn(ctx, authService, sessionService)
		ctx.Next()
	}
}

func ensureLoggedIn(ctx *gin.Context, authService *service.AuthService, sessionService *service.SessionService) {
	// accessCookieString, err := ctx.Cookie(repository.ACCESS_TOKEN.Type)
	// check if the access cookie is there. If not it means there is also no AccessToken because it is inside the cookie
	// IF access cookie NOT present: check refresh cookie.
	// 		IF refresh cookie NOT present: user is NOT logged in => revoke the refresh token
	// 		IF refresh cookie IS present: user is NOT logged in => check refresh cookie and login user

	// IF access cookie IS present: still check expiry time and if not expired, then we are logged in; if expired then not logged
	// 		=> DON'T invoke a new access Token. Cookie has the same TTL as the token. Therefore token most likely stolen or comprimised

}
