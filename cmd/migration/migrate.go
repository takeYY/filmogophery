package main

import (
	"log"

	"filmogophery/internal/config"
	"filmogophery/internal/db"
	"filmogophery/internal/db/migration"
)

func main() {
	// 設定ファイルの読み込み
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// DB 接続
	db.ConnectDB(conf)

	// drop table if exist
	if err := db.WRITER_DB.Migrator().DropTable(
		&migration.User{},
		&migration.Poster{},
		&migration.MovieSeries{},
		&migration.Genre{},
		&migration.Movie{},
		&migration.WatchMedia{},
		&migration.MovieImpression{},
		&migration.MovieWatchRecord{},
	); err != nil {
		log.Fatalf("Error drop tables: %v", err)
	}

	// create table
	if err := db.WRITER_DB.Migrator().CreateTable(
		&migration.User{},
		&migration.Poster{},
		&migration.MovieSeries{},
		&migration.Genre{},
		&migration.Movie{},
		&migration.WatchMedia{},
		&migration.MovieImpression{},
		&migration.MovieWatchRecord{},
	); err != nil {
		log.Fatalf("Error create tables: %v", err)
	}

	// migrate
	if err := db.WRITER_DB.AutoMigrate(
		&migration.User{},
		&migration.Poster{},
		&migration.MovieSeries{},
		&migration.Genre{},
		&migration.Movie{},
		&migration.WatchMedia{},
		&migration.MovieImpression{},
		&migration.MovieWatchRecord{},
	); err != nil {
		log.Fatalf("Error auto migration: %v", err)
	}
}
