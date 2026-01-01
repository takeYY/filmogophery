package main

import (
	"gorm.io/gen"
	"gorm.io/gen/field"

	"filmogophery/internal/pkg/config"
	"filmogophery/internal/pkg/db"
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
	series := g.GenerateModel("series")
	platforms := g.GenerateModel("platforms")

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
		gen.FieldRelate(field.HasOne, "Series", series, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: field.GormTag{
				"foreignKey": []string{"SeriesID"},
				"references": []string{"ID"},
				"default":    []string{"null"},
				"constraint": []string{"OnDelete:SET NULL"},
			},
		}),
	)

	users := g.GenerateModel("users",
		gen.FieldType("password_hash", "constant.PasswordHasher"),
	)
	refreshTokens := g.GenerateModel("refresh_tokens",
		gen.FieldRelate(field.HasOne, "User", users, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"UserID"},
				"references": []string{"ID"},
			},
		}),
	)

	movieGenres := g.GenerateModel("movie_genres")
	watchlist := g.GenerateModel("watchlist",
		gen.FieldRelate(field.HasOne, "User", users, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"UserID"},
				"references": []string{"ID"},
			},
		}),
		gen.FieldRelate(field.HasOne, "Movie", movies, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"MovieID"},
				"references": []string{"ID"},
			},
		}),
	)
	reviews := g.GenerateModel("reviews",
		gen.FieldRelate(field.HasOne, "User", users, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"UserID"},
				"references": []string{"ID"},
			},
		}),
		gen.FieldRelate(field.HasOne, "Movie", movies, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"MovieID"},
				"references": []string{"ID"},
			},
		}),
	)
	watchHistory := g.GenerateModel("watch_history",
		gen.FieldRelate(field.HasOne, "User", users, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"UserID"},
				"references": []string{"ID"},
			},
		}),
		gen.FieldRelate(field.HasOne, "Movie", movies, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"MovieID"},
				"references": []string{"ID"},
			},
		}),
		gen.FieldRelate(field.HasOne, "Platform", platforms, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": []string{"PlatformID"},
				"references": []string{"ID"},
			},
		}),
	)

	g.ApplyBasic(
		genres,
		movies,
		users,
		refreshTokens,
		series,
		movieGenres,
		platforms,
		watchlist,
		reviews,
		watchHistory,
	)

	g.Execute()
}
