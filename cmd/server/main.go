package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"filmogophery/internal/config"
	"filmogophery/internal/db"
	"filmogophery/internal/health"
)

func main() {
	// 設定ファイルの読み込み
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// DB 接続
	db.ConnectDB(conf)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ハンドラの設定
	healthHandler := health.NewHandler(conf)
	healthHandler.RegisterRoutes(e)

	// サーバーの起動
	serverAddr := ":" + conf.Server.Port
	log.Printf("Starting server on %s", serverAddr)
	e.Logger.Fatal(e.Start(serverAddr))
}
