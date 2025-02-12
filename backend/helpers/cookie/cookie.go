package cookie

import (
	"net/http"
	"time"

	"alexlupatsiy.com/personal-website/backend/repository"
	"github.com/gin-gonic/gin"
)

func SetAccessCookie(ctx *gin.Context, accessToken string, ttl int64) {
	cookie := &http.Cookie{
		Name:     repository.ACCESS_COOKIE.Type,
		Value:    accessToken,
		Expires:  time.Now().Add(time.Duration(ttl) * time.Second),
		MaxAge:   int(ttl),
		Secure:   true,                    // Cookie sent over HTTPS only
		HttpOnly: true,                    // Prevents JavaScript access
		SameSite: http.SameSiteStrictMode, // 🔒 Use Strict, Lax, or None
	}

	http.SetCookie(ctx.Writer, cookie)
}
