package impression

import (
	"context"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
)

type (
	QueryService struct {
		MovieImpressionRepo repositories.IImpressionRepository
	}
)

func NewQueryService(movieImpressionRepo repositories.IImpressionRepository) *QueryService {
	return &QueryService{
		MovieImpressionRepo: movieImpressionRepo,
	}
}

func (qs *QueryService) GetImpressions(ctx context.Context) ([]*model.MovieImpression, error) {
	return qs.MovieImpressionRepo.FindAll(ctx)
}
