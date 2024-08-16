package media

import (
	"context"

	"filmogophery/pkg/gen/model"
)

type (
	QueryService struct {
		WatchMediaRepo IQueryRepository
	}
)

func NewQueryService(watchMediaRepo IQueryRepository) *QueryService {
	return &QueryService{
		WatchMediaRepo: watchMediaRepo,
	}
}

func (qs *QueryService) GetWatchMedia(ctx context.Context) ([]*model.WatchMedia, error) {
	return qs.WatchMediaRepo.Find(ctx)
}
