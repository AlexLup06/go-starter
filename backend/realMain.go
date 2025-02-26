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
	passwordResetDb := db.NewPasswordResetDb()
	sessionsDb := db.NewSessionDb()
	userDb := db.NewUserDb()
	authDb := db.NewAuthDb()

	// services
	tokenService := service.NewTokenService([]byte(cfg.JWTKey))
	mailService := service.NewMailService(cfg.SendGridKey)
	sessionService := service.NewSessionService(sessionsDb, tokenService)
	userService := service.NewUserService(userDb, authDb)
	authService := service.NewAuthService(authDb, userService, tokenService)
	passwordResetService := service.NewPasswordResetService(passwordResetDb, userService, tokenService, mailService)

	// handlers
	router := gin.Default()
	authHandler := handler.NewAuthHandler(router, authService, userService, sessionService, passwordResetService)
	staticHandler := handler.NewStaticHandler(router)
	homeHandler := handler.NewHomeHandler(router)
	privateHandler := handler.NewPrivateHandler(router)

	// static
	staticHandler.Routes(cfg.DevMode)

	// middleware
	gzipMiddleware := gzip.Gzip(gzip.DefaultCompression)
	checkHTMXMiddleware := middleware.CheckHTMXRequest()
	dbHandleMiddleware := middleware.InjectDbHandle(contextDb)
	enureLoggedInMiddleware := middleware.EnsureLoggedIn(sessionService)
	setUserInfoMiddleware := middleware.SetUserInfo(sessionService)

	router.Use(checkHTMXMiddleware)
	router.Use(gzipMiddleware)
	router.Use(setUserInfoMiddleware)

	// Routes
	homeHandler.Routes(dbHandleMiddleware)
	authHandler.Routes(dbHandleMiddleware)
	privateHandler.Routes(enureLoggedInMiddleware, dbHandleMiddleware)

	if err := router.Run(":8080"); err != nil {
		return err
	}

	return nil
}
