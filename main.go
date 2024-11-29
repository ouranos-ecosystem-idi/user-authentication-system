package main

import (
	"fmt"

	"authenticator-backend/config"
	"authenticator-backend/interactor"
	"authenticator-backend/presentation/http/echo/middleware"
	"authenticator-backend/presentation/http/echo/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()

	cfg, err := config.NewConfig()
	if err != nil {
		e.Logger.Error("config error")

		return
	}

	logger, _ := config.LoggerBuild(cfg)

	defer func() {
		err := logger.Sync()
		if err != nil {
			e.Logger.Error("logger initialization error: %v", err)
		}
	}()

	zap.ReplaceGlobals(logger)

	switch cfg.LogLevel {
	case "debug":
		e.Logger.SetLevel(log.DEBUG)
	case "info":
		e.Logger.SetLevel(log.INFO)
	case "default":
		e.Logger.SetLevel(log.ERROR)
	}

	middleware.NewMiddleware(e)

	conn := config.NewDBConnection(cfg)

	firebaseConfig := config.NewFirebaseConfig(cfg)

	i := interactor.NewInteractor(
		cfg,
		conn,
		firebaseConfig,
	)
	h := i.NewAppHandler()
	authMiddleware := i.NewAuthMiddleware()

	router.SetRouter(e, h, cfg, conn, authMiddleware)

	if cfg.Env == "local" {
		e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Server.Port)))
	}
}
