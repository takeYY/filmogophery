package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"filmogophery/internal/config"
	"filmogophery/internal/db"
	"filmogophery/internal/health"
	"filmogophery/internal/impression"
	"filmogophery/internal/media"
	"filmogophery/internal/movie"
	"filmogophery/internal/record"
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

	// CORSの設定追加
	e.Use(middleware.CORS())

	// Router 追加
	newRouter(e, conf)

	// ハンドラの設定
	healthHandler := health.NewHandler(conf)
	healthHandler.RegisterRoutes(e)

	tmdbHandler := tmdb.NewHandler(conf)
	tmdbHandler.RegisterRoutes(e)

	// サーバーの起動
	serverAddr := ":" + conf.Server.Port
	logger.Info().Msgf("Starting server on %s", serverAddr)
	e.Logger.Fatal(e.Start(serverAddr))
}

func newRouter(e *echo.Echo, conf *config.Config) {
	// repository の初期化
	impressionRepo := impression.NewQueryRepository()
	mediaRepo := media.NewQueryRepository()
	recordRepo := record.NewQueryRepository()
	movieQueryRepo := movie.NewQueryRepository()
	movieCommandRepo := movie.NewCommandRepository()

	// サービスの初期化
	impressionQueryService := impression.NewQueryService(*impressionRepo)
	mediaQueryService := media.NewQueryService(*mediaRepo)
	recordQueryService := record.NewQueryService(*recordRepo)
	movieQueryService := movie.NewQueryService(conf, *movieQueryRepo)
	movieCommandService := movie.NewCommandService(*movieCommandRepo)

	// ハンドラの追加
	impressionHandler := impression.NewHandler(impressionQueryService)
	impressionHandler.RegisterRoutes(e)

	mediaHandler := media.NewHandler(mediaQueryService)
	mediaHandler.RegisterRoutes(e)

	recordHandler := record.NewHandler(recordQueryService)
	recordHandler.RegisterRoutes(e)

	movieHandler := movie.NewHandler(movieQueryService, movieCommandService)
	movieHandler.RegisterRoutes(e)
}
