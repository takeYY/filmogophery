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
	Tmdb struct {
		ACCESS_TOKEN string
	}

	Config struct {
		Server         server
		Logger         Logger
		ReaderDatabase Database
		WriterDatabase Database
		Tmdb           Tmdb
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
		Tmdb: Tmdb{
			ACCESS_TOKEN: os.Getenv("TMDB_ACCESS_TOKEN"),
		},
	}

	return conf
}
