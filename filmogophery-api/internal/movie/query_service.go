package movie

import (
	"context"

	"filmogophery/internal/db"
	"filmogophery/pkg/gen/model"
)

type (
	QueryService struct {
		MovieRepo IQueryRepository
	}
)

func NewQueryService() *QueryService {
	var movieRepo IQueryRepository = &MovieRepository{
		DB: db.READER_DB,
	}

	return &QueryService{
		MovieRepo: movieRepo,
	}
}

func (qs *QueryService) GetMovieDetails(ctx context.Context, movieID *int64) (*model.Movie, error) {
	id := int32(*movieID)
	return qs.MovieRepo.FindByID(ctx, &id)
}

func (qs *QueryService) GetMovies(ctx context.Context) ([]*model.Movie, error) {
	return qs.MovieRepo.Find(ctx)
}
