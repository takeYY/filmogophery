package movie

import (
	"context"

	"github.com/rs/zerolog"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
)

type (
	GetMovieDetailsUseCase interface {
		Run(ctx context.Context, operator *model.Users, movieID int32) (*types.MovieDetail, error)
	}

	getMovieDetailInteractor struct {
		movieService  services.IMovieService
		reviewService services.IReviewService
		tmdbService   services.ITmdbService
	}

	// 並列処理用のチャンネル
	tmdbResult struct {
		data *types.TmdbMovieDetail
		err  error
	}
	reviewResult struct {
		data *model.Reviews
		err  error
	}
)

func NewGetMovieDetailInteractor(
	movieService services.IMovieService,
	reviewService services.IReviewService,
	tmdbService services.ITmdbService,
) GetMovieDetailsUseCase {
	return &getMovieDetailInteractor{
		movieService,
		reviewService,
		tmdbService,
	}
}

func (i *getMovieDetailInteractor) Run(
	ctx context.Context, operator *model.Users, movieID int32,
) (*types.MovieDetail, error) {
	log := zerolog.Ctx(ctx)

	// 映画詳細を取得
	movie, err := i.movieService.GetMovieByID(ctx, movieID)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("successfully get a movie")

	tmdbCh := make(chan tmdbResult, 1)
	reviewCh := make(chan reviewResult, 1)

	go func() {
		tmdb, err := i.tmdbService.GetMovieDetailByID(movie.TmdbID)
		tmdbCh <- tmdbResult{tmdb, err}
	}()

	go func() {
		review, err := i.reviewService.GetReviewByMovieID(ctx, operator, movie)
		reviewCh <- reviewResult{review, err}
	}()

	// 結果を待機
	tmdbRes := <-tmdbCh
	reviewRes := <-reviewCh

	if tmdbRes.err != nil {
		return nil, tmdbRes.err
	}
	log.Debug().Msg("successfully get a movie from tmdb")

	if reviewRes.err != nil {
		return nil, reviewRes.err
	}
	log.Debug().Msg("successfully get a review")

	// 上映時間を更新
	if movie.RuntimeMinutes == constant.DEFAULT_RUNTIME_MINUTES {
		movie.RuntimeMinutes = int32(tmdbRes.data.Runtime)
		err = i.movieService.UpdateRuntimeMinutes(ctx, nil, movie)
		if err != nil {
			return nil, err
		}
	}

	response := &types.MovieDetail{
		Movie: types.Movie{
			ID:             movie.ID,
			Title:          movie.Title,
			Overview:       movie.Overview,
			ReleaseDate:    types.ConvertTime2Date(movie.ReleaseDate),
			RuntimeMinutes: movie.RuntimeMinutes,
			PosterURL:      movie.PosterURL,
			TmdbID:         movie.TmdbID,
			Genres:         types.NewGenresByModel(movie.Genres),
		},
		VoteAverage: tmdbRes.data.GetVoteAverage(),
		VoteCount:   tmdbRes.data.VoteCount,
		Series:      types.NewSeriesByModel(movie.Series),
		Review:      types.NewReviewByModel(reviewRes.data),
	}
	return response, nil
}
