package backend

import (
	"context"
	"fmt"

	"alexlupatsiy.com/personal-website/backend/config"
	"alexlupatsiy.com/personal-website/backend/db"
	"alexlupatsiy.com/personal-website/backend/handler"
	"alexlupatsiy.com/personal-website/backend/service"
	"github.com/gin-contrib/gzip"

	"github.com/sethvargo/go-envconfig"

	"alexlupatsiy.com/personal-website/backend/middleware"
	"github.com/gin-gonic/gin"
)

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

	// storage
	sessionsDb := db.NewSessionDb()
	userDb := db.NewUserDb()
	authDb := db.NewAuthDb()

	// services
	sessionService := service.NewSessionService(sessionsDb, []byte(cfg.JWTKey))
	userService := service.NewUserService(userDb, authDb)
	authService := service.NewAuthService(authDb, userService, sessionService)

	// handlers
	router := gin.Default()
	authHandler := handler.NewAuthHandler(router, authService, userService, sessionService)
	staticHandler := handler.NewStaticHandler(router)
	homeHandler := handler.NewHomeHandler(router)

	// static
	staticHandler.Routes(cfg.DevMode)

	// middleware
	router.Use(middleware.CheckHTMXRequest())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	dbHandleMiddleware := middleware.InjectDbHandle(contextDb)

	// Routes
	homeHandler.Routes(dbHandleMiddleware)
	authHandler.Routes(dbHandleMiddleware)

	if err := router.Run(":8080"); err != nil {
		return err
	}

	return nil
}
