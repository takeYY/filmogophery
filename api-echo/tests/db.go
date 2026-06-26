package tests

import (
	"os"

	"gorm.io/gorm"

	"filmogophery/internal/pkg/config"
	"filmogophery/internal/pkg/db"
)

// SetupTestDB はテスト用DBへの接続を確立する。
// テスト用の *config.Config を直接構築することで os.Setenv によるグローバル状態汚染を回避する。
func SetupTestDB() *gorm.DB {
	conf := &config.Config{
		ReaderDatabase: config.Database{
			DBCore:   os.Getenv("READER_DB_CORE_COUNT"),
			Addr:     os.Getenv("READER_DB_HOST"),
			User:     os.Getenv("READER_DB_USER"),
			Password: os.Getenv("READER_DB_PWD"),
			Name:     "db4test",
		},
		WriterDatabase: config.Database{
			DBCore:   os.Getenv("WRITER_DB_CORE_COUNT"),
			Addr:     os.Getenv("WRITER_DB_HOST"),
			User:     os.Getenv("WRITER_DB_USER"),
			Password: os.Getenv("WRITER_DB_PWD"),
			Name:     "db4test",
		},
	}
	return db.ConnectDB(conf)
}
