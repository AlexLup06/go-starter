package handler

import (
	"alexlupatsiy.com/personal-website/backend/middleware"
	"github.com/gin-gonic/gin"
)

type StaticHandler struct{ router *gin.Engine }

func NewStaticHandler(router *gin.Engine) *StaticHandler {
	return &StaticHandler{router: router}
}

func (s *StaticHandler) Routes(devMode bool) {
	var staticBasePath string
	if !devMode {
		staticBasePath = "/root/public"
	} else {
		staticBasePath = "./frontend/public"
	}

	static := s.router.Group("/", middleware.ServeGzippedFiles(!devMode))
	{
		static.GET("/js/*filepath", middleware.ServeStaticFiles("./frontend/src/js"))
		static.GET("/css/*filepath", middleware.ServeStaticFiles("./frontend/src/css"))
		static.GET("/public/*filepath", middleware.ServeStaticFiles(staticBasePath))
	}
}
