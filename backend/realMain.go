package backend

import (
	"context"
	"fmt"

	"alexlupatsiy.com/personal-website/backend/config"
	"alexlupatsiy.com/personal-website/backend/db"

	"github.com/sethvargo/go-envconfig"

	"alexlupatsiy.com/personal-website/backend/middleware"
	"alexlupatsiy.com/personal-website/frontend/src/views"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

func RealMain() error {

	cfg := config.Config{}
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return fmt.Errorf("can't inject env variables to dbCfg: %w", err)
	}

	if !cfg.DevMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init db
	dbClient, err := db.NewClient(cfg.DbConfig)
	if err != nil {
		return fmt.Errorf("can't create new db: %w", err)
	}
	contextDb := db.NewContextDb(dbClient.GormDb())
	fmt.Println(contextDb)

	var staticBasePath string
	if !cfg.DevMode {
		staticBasePath = "/root/public"
	} else {
		staticBasePath = "./frontend/public"
	}

	router := gin.Default()

	static := router.Group("/", middleware.ServeGzippedFiles(!cfg.DevMode))
	{
		static.GET("/js/*filepath", middleware.ServeStaticFiles("./frontend/src/js"))
		static.GET("/css/*filepath", middleware.ServeStaticFiles("./frontend/src/css"))
		static.GET("/public/*filepath", middleware.ServeStaticFiles(staticBasePath))
	}

	router.Use(middleware.CheckHTMXRequest())

	router.GET("/", func(c *gin.Context) {
		render(c, 200, views.Home())
	})

	if err := router.Run(":8080"); err != nil {
		return err
	}

	return nil
}
