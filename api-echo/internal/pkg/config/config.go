package config

import "os"

// アプリケーションの設定を格納する構造体
type (
	server struct {
		Port string
	}
	Logger struct {
		Level  string
		Format string
	}
	Database struct {
		DBCore   string
		Addr     string
		User     string
		Password string
		Name     string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       string
	}
	Tmdb struct {
		ACCESS_TOKEN string
	}
	Token struct {
		JWT_SECRET string
	}

	Config struct {
		Server         server
		Logger         Logger
		ReaderDatabase Database
		WriterDatabase Database
		Redis          Redis
		Tmdb           Tmdb
		Token          Token
	}
)

var (
	conf *Config = nil
)

// 設定ファイルを読み込み、Config構造体を返す
func LoadConfig() *Config {
	if conf != nil {
		return conf
	}

	conf = &Config{
		Server: server{
			Port: os.Getenv("SERVER_PORT"),
		},
		Logger: Logger{
			Level:  os.Getenv("LOG_LEVEL"),
			Format: os.Getenv("LOG_FORMAT"),
		},
		ReaderDatabase: Database{
			DBCore:   os.Getenv("READER_DB_CORE_COUNT"),
			Addr:     os.Getenv("READER_DB_HOST"),
			User:     os.Getenv("READER_DB_USER"),
			Password: os.Getenv("READER_DB_PWD"),
			Name:     os.Getenv("READER_DB_NAME"),
		},
		WriterDatabase: Database{
			DBCore:   os.Getenv("WRITER_DB_CORE_COUNT"),
			Addr:     os.Getenv("WRITER_DB_HOST"),
			User:     os.Getenv("WRITER_DB_USER"),
			Password: os.Getenv("WRITER_DB_PWD"),
			Name:     os.Getenv("WRITER_DB_NAME"),
		},
		Redis: Redis{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       os.Getenv("REDIS_DB"),
		},
		Tmdb: Tmdb{
			ACCESS_TOKEN: os.Getenv("TMDB_ACCESS_TOKEN"),
		},
		Token: Token{
			JWT_SECRET: os.Getenv("JWT_SECRET"),
		},
	}

	return conf
}
