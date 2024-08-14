package media

import (
	"context"

	"filmogophery/internal/db"
	"filmogophery/pkg/gen/model"
)

type (
	QueryService struct {
		WatchMediaRepo IQueryRepository
	}
)

func NewQueryService() *QueryService {
	var watchMediaRepo IQueryRepository = &WatchMediaRepository{
		DB: db.READER_DB,
	}

	return &QueryService{
		WatchMediaRepo: watchMediaRepo,
	}
}

func (qs *QueryService) GetWatchMedia(ctx context.Context) ([]*model.WatchMedia, error) {
	return qs.WatchMediaRepo.Find(ctx)
}
