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
	"filmogophery/internal/health"
	"filmogophery/internal/impression"
	"filmogophery/internal/media"
	"filmogophery/internal/movie"
	"filmogophery/internal/pkg/logger"
	"filmogophery/internal/pkg/tokenizer"
	"filmogophery/internal/record"
	"filmogophery/internal/tmdb"
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

	api.RegisterV0Routes(e)

	return e
}

func startServer(e *echo.Echo, conf *config.Config, logger zerolog.Logger) {
	serverAddr := ":" + conf.Server.Port
	logger.Info().Msgf("Starting server on %s", serverAddr)
	e.Logger.Fatal(e.Start(serverAddr))
}

func Old() {
	// 設定ファイルの読み込み
	conf := config.LoadConfig()

	// ロガーの初期化と取得
	logger.InitializeLogger(&conf.Logger)
	logger := logger.GetLogger()

	// DB 接続
	gormDB := db.ConnectDB(conf)

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
	newRouter(e, conf, gormDB)
	api.RegisterV0Routes(e)

	// ハンドラの設定
	healthHandler := health.NewHandler(conf)
	healthHandler.RegisterRoutes(e)

	// サーバーの起動
	serverAddr := ":" + conf.Server.Port
	logger.Info().Msgf("Starting server on %s", serverAddr)
	e.Logger.Fatal(e.Start(serverAddr))
}

func newRouter(e *echo.Echo, conf *config.Config, gormDB *gorm.DB) {
	// ----- Init Repository ----- //
	genreRepo := repositories.NewGenreRepository(gormDB)
	impressionRepo := repositories.NewImpressionRepository(gormDB)
	mediaRepo := repositories.NewMediaRepository(gormDB)
	recordRepo := repositories.NewRecordRepository(gormDB)
	movieRepo := repositories.NewMovieRepository(gormDB)

	// 外部 API
	tmdbClient := tmdb.NewTmdbClient(conf)

	// ----- サービスの初期化 ----- //
	// impression
	impressionQueryService := impression.NewQueryService(impressionRepo)
	// media
	mediaQueryService := media.NewQueryService(mediaRepo)
	// record
	recordQueryService := record.NewQueryService(recordRepo)
	recordCommandService := record.NewCommandService(recordRepo, mediaRepo, impressionRepo)
	// movie
	movieQueryService := movie.NewQueryService(conf, movieRepo, recordRepo)
	movieCommandService := movie.NewCommandService(movieRepo, genreRepo, impressionRepo, *tmdbClient)
	// 外部 API
	tmdbService := tmdb.NewTmdbService(*tmdbClient)

	// ハンドラの追加
	impressionHandler := impression.NewHandler(impressionQueryService)
	impressionHandler.RegisterRoutes(e)

	mediaHandler := media.NewHandler(mediaQueryService)
	mediaHandler.RegisterRoutes(e)

	recordHandler := record.NewHandler(recordQueryService, recordCommandService)
	recordHandler.RegisterRoutes(e)

	movieHandler := movie.NewHandler(movieQueryService, movieCommandService)
	movieHandler.RegisterRoutes(e)

	tmdbHandler := tmdb.NewHandler(tmdbService)
	tmdbHandler.RegisterRoutes(e)
}
