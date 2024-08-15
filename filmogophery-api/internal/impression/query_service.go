package impression

import (
	"context"

	"filmogophery/internal/db"
	"filmogophery/pkg/gen/model"
)

type (
	QueryService struct {
		MovieImpressionRepo IQueryRepository
	}
)

func NewQueryService() *QueryService {
	var movieImpressionRepo IQueryRepository = &MovieImpressionRepository{
		DB: db.READER_DB,
	}

	return &QueryService{
		MovieImpressionRepo: movieImpressionRepo,
	}
}

func (qs *QueryService) GetImpressions(ctx context.Context) ([]*model.MovieImpression, error) {
	return qs.MovieImpressionRepo.Find(ctx)
}
