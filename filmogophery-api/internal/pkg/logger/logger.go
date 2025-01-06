package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"filmogophery/internal/config"
)

// zerolog のロガーを初期化
func InitializeLogger(conf *config.Logger) {
	// 標準出力
	if conf.Format == "json" {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "2006-01-02 15:04:05.000000",
		}).With().Timestamp().Logger()
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// ログレベルの設定
	level, err := zerolog.ParseLevel(conf.Level)
	if err != nil {
		log.Logger.Error().Msgf("Invalid log level %s, defaulting to info", conf.Level)
		level = zerolog.InfoLevel
	}
	log.Logger = log.Logger.Level(level)
}

// zerolog のロガーインスタンスを取得
func GetLogger() zerolog.Logger {
	return log.Logger
}
