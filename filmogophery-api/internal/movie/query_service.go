package movie

import (
	"context"

	"filmogophery/internal/config"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
	"filmogophery/internal/record"
	"filmogophery/internal/tmdb"
)

type (
	QueryService struct {
		MovieRepo  IQueryRepository
		RecordRepo record.IQueryRepository
		TmdbClient tmdb.ITmdbClient
	}
)

func NewQueryService(conf *config.Config, movieRepo IQueryRepository) *QueryService {
	return &QueryService{
		MovieRepo:  movieRepo,
		RecordRepo: *record.NewQueryRepository(),
		TmdbClient: *tmdb.NewTmdbClient(conf),
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

	var genres []*string = make([]*string, 0)
	for _, m := range movie.Genres {
		genres = append(genres, m.Name)
	}

	tmdbMovie, err := qs.TmdbClient.GetMovieDetail(movie.TmdbID)
	if err != nil {
		logger.Info().Msg("failed to get a detail movie from tmdb")
		return nil, err
	}
	logger.Info().Msgf("success to get a detail movie from tmdb")

	var records []*model.MovieWatchRecord = make([]*model.MovieWatchRecord, 0)
	if movie.MovieImpression != nil {
		records, err = qs.RecordRepo.FindByImpressionID(ctx, &movie.MovieImpression.ID)
		if err != nil {
			logger.Info().Msg("failed to get watch records")
			return nil, err
		}
		logger.Info().Msgf("success to get watch records")
	}

	var watchRecords []*WatchRecordDetailDto = make([]*WatchRecordDetailDto, 0)
	for _, r := range records {
		watchRecord := &WatchRecordDetailDto{
			WatchDate:  r.WatchDate.Format("2006-01-02"),
			WatchMedia: *r.WatchMedia.Name,
		}
		watchRecords = append(watchRecords, watchRecord)
	}

	var series SeriesDetailDto
	if movie.Series != nil {
		series = SeriesDetailDto{
			Name:      movie.Series.Name,
			PosterURL: *movie.Series.PosterURL,
		}
	}
	var impression ImpressionDetailDto
	if movie.MovieImpression != nil {
		impression = ImpressionDetailDto{
			ID:     movie.MovieImpression.ID,
			Status: movie.MovieImpression.Status,
			Rating: movie.MovieImpression.Rating,
			Note:   movie.MovieImpression.Note,
		}
	}

	movieDetail := &MovieDetailDto{
		ID:           movie.ID,
		Title:        movie.Title,
		Overview:     movie.Overview,
		ReleaseDate:  movie.ReleaseDate.Format("2006-01-02"),
		RunTime:      movie.RunTime,
		Genres:       genres,
		PosterURL:    *movie.PosterURL,
		VoteAverage:  tmdbMovie.VoteAverage,
		VoteCount:    tmdbMovie.VoteCount,
		Series:       &series,
		Impression:   &impression,
		WatchRecords: watchRecords,
	}

	return movieDetail, nil
}

func (qs *QueryService) GetMovies(ctx context.Context) ([]*model.Movie, error) {
	return qs.MovieRepo.Find(ctx)
}
