package impression

import (
	"context"

	"filmogophery/internal/pkg/gen/model"
)

type (
	QueryService struct {
		MovieImpressionRepo IQueryRepository
	}
)

func NewQueryService(movieImpressionRepo IQueryRepository) *QueryService {
	return &QueryService{
		MovieImpressionRepo: movieImpressionRepo,
	}
}

func (qs *QueryService) GetImpressions(ctx context.Context) ([]*model.MovieImpression, error) {
	return qs.MovieImpressionRepo.Find(ctx)
}
