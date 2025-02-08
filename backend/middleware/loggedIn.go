package middleware

import (
	"alexlupatsiy.com/personal-website/backend/service"
	"github.com/gin-gonic/gin"
)

func EnsureLoggedIn(authService *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ensureLoggedIn(ctx, authService)
		ctx.Next()
	}
}

func ensureLoggedIn(ctx *gin.Context, authService *service.AuthService) {

}
