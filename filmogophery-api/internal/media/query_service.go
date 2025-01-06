package media

import (
	"context"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
)

type (
	QueryService struct {
		WatchMediaRepo repositories.IMediaRepository
	}
)

func NewQueryService(watchMediaRepo repositories.IMediaRepository) *QueryService {
	return &QueryService{
		WatchMediaRepo: watchMediaRepo,
	}
}

func (qs *QueryService) GetWatchMedia(ctx context.Context) ([]*model.WatchMedia, error) {
	return qs.WatchMediaRepo.FindAll(ctx)
}
