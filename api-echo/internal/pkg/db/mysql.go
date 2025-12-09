package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	defaultMySQL "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/config"
)

var (
	gormDB *gorm.DB = nil
)

func ConnectDB(conf *config.Config) *gorm.DB {
	if gormDB != nil {
		return gormDB
	}

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Error(fmt.Printf("failed to load location %s", err))
	}

	readerDB := defaultMySQL.Config{
		DBName:               conf.ReaderDatabase.Name,
		User:                 conf.ReaderDatabase.User,
		Passwd:               conf.ReaderDatabase.Password,
		Addr:                 conf.ReaderDatabase.Addr,
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
	writerDB := defaultMySQL.Config{
		DBName:               conf.WriterDatabase.Name,
		User:                 conf.WriterDatabase.User,
		Passwd:               conf.WriterDatabase.Password,
		Addr:                 conf.WriterDatabase.Addr,
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}

	db, err := sql.Open("mysql", writerDB.FormatDSN())
	if err != nil {
		log.Error(fmt.Printf("failed to connect mysql %s", err))
	}

	// Connection Config
	coreCount, err := strconv.Atoi(conf.WriterDatabase.DBCore)
	if err != nil {
		log.Error(fmt.Printf("failed to convert from string to integer: %s", conf.WriterDatabase.DBCore))
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(coreCount * 2)
	db.SetConnMaxLifetime(time.Hour)

	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		CreateBatchSize: 1000,
		TranslateError:  true,
	})
	if err != nil {
		log.Error(fmt.Printf("failed to use gorm %s", err))
	}

	// DbResolver Config
	gormDB.Use(
		dbresolver.Register(
			// Read Replica Config
			dbresolver.Config{
				Replicas: []gorm.Dialector{mysql.Open(readerDB.FormatDSN())},
			},
		),
	)

	return gormDB
}
