package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"filmogophery/internal/app/api"
	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/services"
	"filmogophery/internal/config"
	"filmogophery/internal/db"
	"filmogophery/internal/pkg/logger"
)

func main() {
	app := fx.New(
		fx.Provide( // Requirements
			config.LoadConfig, // 設定ファイル
			newLogger,         // ロガー
			db.ConnectDB,      // DB
		),
		fx.Provide( // Repositories
			repositories.NewGenreRepository,
			repositories.NewImpressionRepository,
			repositories.NewMediaRepository,
			repositories.NewMovieRepository,
			repositories.NewRecordRepository,
		),
		fx.Provide( // Services
			services.NewMovieService,
		),
		fx.Provide(
			services.NewServiceContainer,
		),
		fx.Provide(
			newEchoServer,
		),
		api.RegisterV0Routes(),
		api.RegisterV1Routes(),
		fx.Invoke(startServer),
	)
	app.Run()
}

func newLogger(conf *config.Config) zerolog.Logger {
	logger.InitializeLogger(&conf.Logger)
	return logger.GetLogger()
}

func newEchoServer(conf *config.Config, gormDB *gorm.DB, serviceContainer services.IServiceContainer) *echo.Echo {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}

func startServer(e *echo.Echo, conf *config.Config, logger zerolog.Logger) {
	serverAddr := ":" + conf.Server.Port
	logger.Info().Msgf("Starting server on %s", serverAddr)
	e.Logger.Fatal(e.Start(serverAddr))
}
