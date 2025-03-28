package handler

import (
	"alexlupatsiy.com/personal-website/backend/helpers/render"
	"alexlupatsiy.com/personal-website/frontend/src/views"
	"github.com/gin-gonic/gin"
)

type HomeHandler struct{ router *gin.Engine }

func NewHomeHandler(router *gin.Engine) *HomeHandler {
	return &HomeHandler{router: router}
}

func (h *HomeHandler) Routes(dbHandleMiddleware gin.HandlerFunc) {
	h.router.GET("/", h.home)
	h.router.GET("/test", h.test)
}

func (h *HomeHandler) home(ctx *gin.Context) {
	render.Render(ctx, 200, views.Home())
}

func (h HomeHandler) test(ctx *gin.Context) {
	render.Render(ctx, 200, views.Test())
}
