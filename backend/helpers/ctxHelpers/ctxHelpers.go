package ctxHelpers

import (
	"context"

	"github.com/gin-gonic/gin"
)

func SetKV(ctx *gin.Context, key, value string) {
	newCtx := context.WithValue(ctx.Request.Context(), key, value)
	ctx.Request = ctx.Request.WithContext(newCtx)
}
