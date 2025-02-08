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

func (h *HomeHandler) Routes() {
	h.router.GET("/", h.home)
}

func (h *HomeHandler) home(ctx *gin.Context) {
	render.Render(ctx, 200, views.Home())
}
