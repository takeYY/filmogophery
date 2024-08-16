package record

import (
	"context"

	"filmogophery/pkg/gen/model"
)

type (
	QueryService struct {
		MovieWatchRecordRepo IQueryRepository
	}
)

func NewQueryService(movieWatchRecordRepo IQueryRepository) *QueryService {
	return &QueryService{
		MovieWatchRecordRepo: movieWatchRecordRepo,
	}
}

func (qs *QueryService) GetWatchRecords(ctx context.Context) ([]*model.MovieWatchRecord, error) {
	return qs.MovieWatchRecordRepo.Find(ctx)
}
