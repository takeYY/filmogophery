package config

import "os"

// アプリケーションの設定を格納する構造体
type (
	server struct {
		Port string
	}
	Database struct {
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
		ReaderDatabase Database
		WriterDatabase Database
		Tmdb           Tmdb
	}
)

// 設定ファイルを読み込み、Config構造体を返す
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: server{
			Port: os.Getenv("SERVER_PORT"),
		},
		ReaderDatabase: Database{
			Addr:     os.Getenv("READER_DB_HOST"),
			User:     os.Getenv("READER_DB_USER"),
			Password: os.Getenv("READER_DB_PWD"),
			Name:     os.Getenv("READER_DB_NAME"),
		},
		WriterDatabase: Database{
			Addr:     os.Getenv("WRITER_DB_HOST"),
			User:     os.Getenv("WRITER_DB_USER"),
			Password: os.Getenv("WRITER_DB_PWD"),
			Name:     os.Getenv("WRITER_DB_NAME"),
		},
		Tmdb: Tmdb{
			ACCESS_TOKEN: os.Getenv("TMDB_ACCESS_TOKEN"),
		},
	}

	return config, nil
}
