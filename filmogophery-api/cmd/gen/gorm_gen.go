package main

import (
	"gorm.io/gen"
	"gorm.io/gen/field"

	"filmogophery/internal/db"
	"filmogophery/internal/pkg/config"
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
	conf := config.LoadConfig()
	gormDB := db.ConnectDB(conf)

	g.UseDB(gormDB)

	// 外部キーの関係を持つ構造体を作る
	genres := g.GenerateModel("genres",
		gen.FieldRelate(field.Many2Many, "Movies", g.GenerateModel("movies"), &field.RelateConfig{
			RelateSlicePointer: true,
			GORMTag:            field.GormTag{"many2many": []string{"movie_genres"}},
		}))
	movies := g.GenerateModel("movies",
		gen.FieldRelate(field.Many2Many, "Genres", genres, &field.RelateConfig{
			RelateSlicePointer: true,
			GORMTag: field.GormTag{
				"many2many":      []string{"movie_genres"},
				"foreignKey":     []string{"ID"},
				"joinForeignKey": []string{"movie_id"},
				"references":     []string{"ID"},
				"joinReferences": []string{"genre_id"},
			},
		}),
		gen.FieldRelate(field.HasOne, "Series", g.GenerateModel("series"), &field.RelateConfig{
			RelatePointer: true,
			GORMTag: field.GormTag{
				"foreignKey": []string{"SeriesID"},
				"references": []string{"ID"},
				"default":    []string{"null"},
				"constraint": []string{"OnDelete:SET NULL"},
			},
		}),
	)
	users := g.GenerateModel("users")
	series := g.GenerateModel("series")
	movieGenres := g.GenerateModel("movie_genres")
	platforms := g.GenerateModel("platforms")
	watchlist := g.GenerateModel("watchlist")
	reviews := g.GenerateModel("reviews")
	watchHistory := g.GenerateModel("watch_history")

	g.ApplyBasic(
		genres,
		movies,
		users,
		series,
		movieGenres,
		platforms,
		watchlist,
		reviews,
		watchHistory,
	)

	g.Execute()
}
