package movie

import (
	"context"

	"filmogophery/internal/config"
	"filmogophery/internal/db"
	"filmogophery/internal/record"
	"filmogophery/internal/tmdb"
	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/logger"
)

type (
	QueryService struct {
		MovieRepo  IQueryRepository
		RecordRepo record.IQueryRepository
		TmdbClient tmdb.ITmdbClient
	}
)

func NewQueryService(conf *config.Config) *QueryService {
	var movieRepo IQueryRepository = &MovieRepository{
		DB: db.READER_DB,
	}

	return &QueryService{
		MovieRepo:  movieRepo,
		RecordRepo: record.NewQueryRepository(),
		TmdbClient: tmdb.NewTmdbClient(conf),
	}
}

func (qs *QueryService) GetMovieDetails(ctx context.Context, movieID *int64) (*MovieDetailDto, error) {
	logger := logger.GetLogger()

	id := int32(*movieID)
	movie, err := qs.MovieRepo.FindByID(ctx, &id)
	if err != nil {
		logger.Info().Msg("failed to get a movie")
		return nil, err
	}
	logger.Info().Msg("success to get a movie")

	var genres []string
	for _, m := range movie.Genres {
		genres = append(genres, *m.Name)
	}

	tmdbMovie, err := qs.TmdbClient.GetMovieDetail(movie.TmdbID)
	if err != nil {
		logger.Info().Msg("failed to get a detail movie from tmdb")
		return nil, err
	}
	logger.Info().Msgf("success to get a detail movie from tmdb")

	records, err := qs.RecordRepo.FindByImpressionID(ctx, &movie.MovieImpression.ID)
	if err != nil {
		logger.Info().Msg("failed to get watch records")
		return nil, err
	}
	logger.Info().Msgf("success to get watch records")

	var watchRecords []*WatchRecordDetailDto
	for _, r := range records {
		watchRecord := &WatchRecordDetailDto{
			WatchDate:  r.WatchDate.Format("2006-01-02"),
			WatchMedia: *r.WatchMedia.Name,
		}
		watchRecords = append(watchRecords, watchRecord)
	}

	movieDetail := &MovieDetailDto{
		ID:          movie.ID,
		Title:       movie.Title,
		Overview:    movie.Overview,
		ReleaseDate: movie.ReleaseDate.Format("2006-01-02"),
		RunTime:     movie.RunTime,
		Genres:      genres,
		PosterURL:   movie.Poster.URL,
		VoteAverage: tmdbMovie.VoteAverage,
		VoteCount:   tmdbMovie.VoteCount,
		Series: &SeriesDetailDto{
			Name:      movie.Series.Name,
			PosterURL: movie.Series.Poster.URL,
		},
		Impression: &ImpressionDetailDto{
			ID:     movie.MovieImpression.ID,
			Status: movie.MovieImpression.Status,
			Rating: movie.MovieImpression.Rating,
			Note:   movie.MovieImpression.Note,
		},
		WatchRecords: watchRecords,
	}

	return movieDetail, nil
}

func (qs *QueryService) GetMovies(ctx context.Context) ([]*model.Movie, error) {
	return qs.MovieRepo.Find(ctx)
}
