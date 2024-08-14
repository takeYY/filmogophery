package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"filmogophery/internal/config"
	"filmogophery/internal/db"
	"filmogophery/internal/health"
	"filmogophery/internal/movie"
	"filmogophery/internal/tmdb"
	"filmogophery/pkg/logger"
	"filmogophery/pkg/tokenizer"
)

func main() {
	// ロガーの初期化と取得
	logger.InitializeLogger("info")
	logger := logger.GetLogger()

	// 設定ファイルの読み込み
	conf, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Msgf("Error loading config: %v", err)
	}

	// DB 接続
	db.ConnectDB(conf)

	// 形態素解析器の初期化
	er := tokenizer.NewTokenizer()
	if er != nil {
		logger.Fatal().Msgf("failed to create tokenizer: %v", er)
	}

	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			return nil
		},
	}))
	e.Use(middleware.Recover())

	// ハンドラの設定
	healthHandler := health.NewHandler(conf)
	healthHandler.RegisterRoutes(e)

	movieHandler := movie.NewHandler()
	movieHandler.RegisterRoutes(e)

	tmdbHandler := tmdb.NewHandler(conf)
	tmdbHandler.RegisterRoutes(e)

	// サーバーの起動
	serverAddr := ":" + conf.Server.Port
	logger.Info().Msgf("Starting server on %s", serverAddr)
	e.Logger.Fatal(e.Start(serverAddr))
}