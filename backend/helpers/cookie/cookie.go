package cookie

import (
	"net/http"
	"time"

	"alexlupatsiy.com/personal-website/backend/repository"
	"github.com/gin-gonic/gin"
)

func SetCookie(ctx *gin.Context, token string, cookieType repository.CookieType, ttl int64) {
	cookie := &http.Cookie{
		Name:     cookieType.Type,
		Value:    token,
		Expires:  time.Now().Add(time.Duration(ttl) * time.Second),
		MaxAge:   int(ttl),
		Secure:   true,                    // Cookie sent over HTTPS only
		HttpOnly: true,                    // Prevents JavaScript access
		SameSite: http.SameSiteStrictMode, // 🔒 Use Strict, Lax, or None
		Path:     "/",
	}
	http.SetCookie(ctx.Writer, cookie)
}

func DeleteCookie(ctx *gin.Context, cookieType repository.CookieType) {
	cookie := &http.Cookie{
		Name:     cookieType.Type,
		Value:    "",
		Expires:  time.Now(),
		MaxAge:   -1,
		Secure:   true,                    // Cookie sent over HTTPS only
		HttpOnly: true,                    // Prevents JavaScript access
		SameSite: http.SameSiteStrictMode, // 🔒 Use Strict, Lax, or None
		Path:     "/",
	}

	http.SetCookie(ctx.Writer, cookie)
}
