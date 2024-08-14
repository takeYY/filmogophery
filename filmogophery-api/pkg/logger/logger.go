package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// zerolog のロガーを初期化
func InitializeLogger(logLevel string) {
	// 標準出力に出力
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// ログレベルの設定
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Logger.Error().Msgf("Invalid log level %s, defaulting to info", logLevel)
		level = zerolog.InfoLevel
	}
	log.Logger = log.Logger.Level(level)
}

// zerolog のロガーインスタンスを取得
func GetLogger() zerolog.Logger {
	return log.Logger
}
