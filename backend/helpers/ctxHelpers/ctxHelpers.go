package ctxHelpers

import (
	"context"

	"github.com/gin-gonic/gin"
)

func SetKV(ctx *gin.Context, key, value string) {
	newCtx := context.WithValue(ctx.Request.Context(), key, value)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

type contextKey string

const (
	loggedInKey     contextKey = "logged_in"
	weekLoggedInKey contextKey = "week_logged_in"
	userIdKey       contextKey = "user_id"
	usernameKey     contextKey = "username"
	emailKey        contextKey = "email"
)

// WithIsLoggedInCtx returns a new context that has a key indicating that the user is logged in.
func WithIsLoggedInCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggedInKey, true)
}

func WithIsWeekLoggedInCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, weekLoggedInKey, true)
}

func WithUserIdCtx(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}

func WithUsernameCtx(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

func WithEmailCtx(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailKey, email)
}

func IsLoggedInCtx(ctx context.Context) bool {
	loggedIn, ok := ctx.Value(loggedInKey).(bool)
	if !ok {
		return false
	}
	return loggedIn
}

func IsWeekLoggedInCtx(ctx context.Context) bool {
	weekLoggedIn, ok := ctx.Value(weekLoggedInKey).(bool)
	if !ok {
		return false
	}
	return weekLoggedIn
}

func UserIdCtx(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(userIdKey).(string)
	return userId, ok
}

func UsernameCtx(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}

func EmailCtx(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(emailKey).(string)
	return email, ok
}
