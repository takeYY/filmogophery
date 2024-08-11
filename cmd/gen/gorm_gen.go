package main

import (
	"log"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"filmogophery/internal/config"
	"filmogophery/internal/db"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./pkg/gen/query",
		Mode:         gen.WithoutContext,
		ModelPkgPath: "./pkg/gen/model",
	})

	// 設定ファイルの読み込み
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	db.ConnectDB(conf)

	g.UseDB(db.READER_DB)

	all := g.GenerateAllTable()

	// 外部キーの関係を持つ構造体を作る
	poster := g.GenerateModel("poster")
	movie_series := g.GenerateModel("movie_series", gen.FieldRelate(field.HasOne, "Poster", poster, &field.RelateConfig{
		GORMTag: field.GormTag{"foreignKey": []string{"ID"}},
	}))

	genre := g.GenerateModel("genre", gen.FieldRelate(field.Many2Many, "Movies", g.GenerateModel("movie"), &field.RelateConfig{
		RelateSlice: true,
		GORMTag:     field.GormTag{"many2many": []string{"movie_genres"}},
	}))
	movie := g.GenerateModel("movie", gen.FieldRelate(field.Many2Many, "Genres", genre, &field.RelateConfig{
		RelateSlice: true,
		GORMTag:     field.GormTag{"many2many": []string{"movie_genre"}},
	}), gen.FieldRelate(field.HasOne, "Poster", poster, &field.RelateConfig{
		GORMTag: field.GormTag{"foreignKey": []string{"ID"}},
	}), gen.FieldRelate(field.HasOne, "Series", movie_series, &field.RelateConfig{
		GORMTag: field.GormTag{"foreignKey": []string{"ID"}},
	}))

	movie_impression := g.GenerateModel("movie_impression", gen.FieldRelate(field.HasOne, "Movie", movie, &field.RelateConfig{
		GORMTag: field.GormTag{"foreignKey": []string{"ID"}},
	}))

	watch_media := g.GenerateModel("watch_media")
	movie_watch_record := g.GenerateModel("movie_watch_record", gen.FieldRelate(field.HasOne, "WatchMedia", watch_media, &field.RelateConfig{
		GORMTag: field.GormTag{"foreignKey": []string{"ID"}},
	}))

	g.ApplyBasic(poster, movie_series, genre, movie, movie_impression, watch_media, movie_watch_record)
	g.ApplyBasic(all...)

	g.Execute()
}
