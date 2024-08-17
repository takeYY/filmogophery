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
		OutPath:        "./internal/pkg/gen/query",
		Mode:           gen.WithoutContext,
		ModelPkgPath:   "./internal/pkg/gen/model",
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
	movie_series := g.GenerateModel("movie_series")

	watch_media := g.GenerateModel("watch_media")
	movie_watch_record := g.GenerateModel("movie_watch_record",
		gen.FieldRelate(field.HasOne, "WatchMedia", watch_media, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"WatchMediaID"},
				"references": []string{"ID"},
			},
		}),
	)

	movie_impression := g.GenerateModel("movie_impression",
		gen.FieldRelate(field.HasOne, "Movie", g.GenerateModel("movie"), &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"MovieID"},
				"references": []string{"ID"},
			},
		}),
		gen.FieldRelate(field.HasMany, "WatchRecords", movie_watch_record, &field.RelateConfig{
			RelateSlicePointer: true,
			GORMTag: field.GormTag{
				"foreignKey": []string{"MovieImpressionID"},
				"references": []string{"ID"},
			},
		}),
	)

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
		gen.FieldRelate(field.HasOne, "Series", movie_series, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: field.GormTag{
				"foreignKey": []string{"SeriesID"},
				"references": []string{"ID"},
				"default":    []string{"null"},
				"constraint": []string{"OnUpdate:SET NULL", "OnDelete:SET NULL"},
			},
		}),
		gen.FieldRelate(field.BelongsTo, "MovieImpression", movie_impression, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: field.GormTag{
				"foreignKey": []string{"MovieID"},
				"references": []string{"ID"},
				// "default":    []string{"null"},
				// "constraint": []string{"OnUpdate:SET NULL", "OnDelete:SET NULL"},
			},
		}),
	)

	g.ApplyBasic(movie_series, genre, movie, movie_impression, watch_media, movie_watch_record)
	g.ApplyBasic(all...)

	g.Execute()
}
