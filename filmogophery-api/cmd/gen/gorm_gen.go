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
		OutPath:        "./pkg/gen/query",
		Mode:           gen.WithoutContext,
		ModelPkgPath:   "./pkg/gen/model",
		FieldNullable:  true,
		FieldCoverable: true,
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
		GORMTag: field.GormTag{
			"foreignKey": []string{"PosterID"},
			"references": []string{"ID"},
			"default":    []string{"null"},
		},
	}))

	genre := g.GenerateModel("genre",
		gen.FieldRelate(field.Many2Many, "Movies", g.GenerateModel("movie"), &field.RelateConfig{
			RelateSlicePointer: true,
			GORMTag:            field.GormTag{"many2many": []string{"movie_genres"}},
		}))
	movie := g.GenerateModel("movie",
		gen.FieldRelate(field.Many2Many, "Genres", genre, &field.RelateConfig{
			RelateSlicePointer: true,
			GORMTag:            field.GormTag{"many2many": []string{"movie_genres"}},
		}),
		// gen.FieldIgnore("PosterID"),
		// gen.FieldIgnore("SeriesID"),
		gen.FieldRelate(field.HasOne, "Poster", poster, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: field.GormTag{
				"foreignKey": []string{"PosterID"},
				"references": []string{"ID"},
				"default":    []string{"null"},
				"constraint": []string{"OnUpdate:SET NULL", "OnDelete:SET NULL"},
			},
		}),
		gen.FieldRelate(field.HasOne, "Series", movie_series, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: field.GormTag{
				"foreignKey": []string{"SeriesID"},
				"references": []string{"ID"},
				"default":    []string{"null"},
				"constraint": []string{"OnUpdate:SET NULL", "OnDelete:SET NULL"},
			},
		}))

	movie_impression := g.GenerateModel("movie_impression", gen.FieldRelate(field.HasOne, "Movie", movie, &field.RelateConfig{
		GORMTag: field.GormTag{
			"foreignKey": []string{"MovieID"},
			"references": []string{"ID"},
			"default":    []string{"null"},
		},
	}))

	watch_media := g.GenerateModel("watch_media")
	movie_watch_record := g.GenerateModel("movie_watch_record", gen.FieldRelate(field.HasOne, "WatchMedia", watch_media, &field.RelateConfig{
		GORMTag: field.GormTag{
			"foreignKey": []string{"WatchMediaID"},
			"references": []string{"ID"},
			"default":    []string{"null"},
		},
	}))

	g.ApplyBasic(poster, movie_series, genre, movie, movie_impression, watch_media, movie_watch_record)
	g.ApplyBasic(all...)

	g.Execute()
}