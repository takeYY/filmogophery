package tests

import (
	"os"

	"gorm.io/gorm"

	"filmogophery/internal/config"
	"filmogophery/internal/db"
)

func SetupTestDB() *gorm.DB {
	// テスト用環境変数を設定
	os.Setenv("READER_DB_NAME", "db4test")
	os.Setenv("WRITER_DB_NAME", "db4test")

	conf := config.LoadConfig()
	return db.ConnectDB(conf)
}
