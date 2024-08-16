package impression

import (
	"time"

	"filmogophery/pkg/gen/model"
)

func InitImpressionNoData() []*model.MovieImpression {
	return make([]*model.MovieImpression, 0)
}

func InitImpression() []*model.MovieImpression {
	rating := float32(4.5)
	note := "テスト感想_1"

	movieOverview := "テスト概要_1"
	moviePosterID := int32(1)
	movieSeriesID := int32(1)
	movieTmdbID := int32(456)
	movies := []model.Movie{
		{
			ID:              1,
			Title:           "テストタイトル_1",
			Overview:        &movieOverview,
			ReleaseDate:     time.Date(2024, 1, 2, 3, 4, 5, 6789, time.Local),
			RunTime:         123,
			PosterID:        &moviePosterID,
			SeriesID:        &movieSeriesID,
			TmdbID:          &movieTmdbID,
			Genres:          nil,
			Poster:          nil,
			Series:          nil,
			MovieImpression: nil,
		},
		{
			ID:              2,
			Title:           "テストタイトル_2",
			Overview:        nil,
			ReleaseDate:     time.Date(2020, 2, 3, 4, 5, 6, 789, time.Local),
			RunTime:         456,
			PosterID:        nil,
			SeriesID:        nil,
			TmdbID:          nil,
			Genres:          nil,
			Poster:          nil,
			Series:          nil,
			MovieImpression: nil,
		},
	}

	watchRecords := []*model.MovieWatchRecord{
		{
			ID:                1,
			MovieImpressionID: 1,
			WatchMediaID:      99,
			WatchDate:         time.Date(2016, 12, 25, 0, 0, 0, 0, time.Local),
			WatchMedia: model.WatchMedia{
				ID:   0,
				Code: "",
				Name: nil,
			},
		},
		{
			ID:                2,
			MovieImpressionID: 1,
			WatchMediaID:      1,
			WatchDate:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			WatchMedia: model.WatchMedia{
				ID:   0,
				Code: "",
				Name: nil,
			},
		},
	}

	impressions := []*model.MovieImpression{
		{
			ID:           1,
			MovieID:      1,
			Status:       true,
			Rating:       &rating,
			Note:         &note,
			Movie:        movies[0],
			WatchRecords: watchRecords,
		},
		{
			ID:           2,
			MovieID:      2,
			Status:       false,
			Rating:       nil,
			Note:         nil,
			Movie:        movies[1],
			WatchRecords: make([]*model.MovieWatchRecord, 0),
		},
	}
	return impressions
}