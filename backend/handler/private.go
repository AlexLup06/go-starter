package handler

import (
	"alexlupatsiy.com/personal-website/backend/helpers/render"
	"alexlupatsiy.com/personal-website/frontend/src/views"
	"github.com/gin-gonic/gin"
)

type PrivateHandler struct{ router *gin.Engine }

func NewPrivateHandler(router *gin.Engine) PrivateHandler {
	return PrivateHandler{router: router}
}

func (p *PrivateHandler) Routes(ensureLoggedInMiddleware gin.HandlerFunc, dbHandleMiddleware gin.HandlerFunc) {
	private := p.router.Group("/private", ensureLoggedInMiddleware, dbHandleMiddleware)
	private.GET("", p.private)
}

func (p *PrivateHandler) private(ctx *gin.Context) {
	render.Render(ctx, 200, views.Private())
}
