package db

import (
	"database/sql"
	"fmt"
	"time"

	default_mysql "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"filmogophery/internal/config"
)

var (
	READER_DB *gorm.DB
	WRITER_DB *gorm.DB
)

func useGORM(conf *config.Database) *gorm.DB {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Error(fmt.Printf("failed to load location %s", err))
	}

	c := default_mysql.Config{
		DBName:               conf.Name,
		User:                 conf.User,
		Passwd:               conf.Password,
		Addr:                 conf.Addr,
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		log.Error(fmt.Printf("failed to connect mysql %s", err))
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Error(fmt.Printf("failed to use gorm %s", err))
	}

	return gormDB
}

func ConnectDB(conf *config.Config) {
	READER_DB = useGORM(&conf.ReaderDatabase)
	WRITER_DB = useGORM(&conf.WriterDatabase)
}
