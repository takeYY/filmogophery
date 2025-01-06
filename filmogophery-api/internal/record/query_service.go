package record

import (
	"context"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
)

type (
	QueryService struct {
		MovieWatchRecordRepo repositories.IRecordRepository
	}
)

func NewQueryService(movieWatchRecordRepo repositories.IRecordRepository) *QueryService {
	return &QueryService{
		MovieWatchRecordRepo: movieWatchRecordRepo,
	}
}

func (qs *QueryService) GetWatchRecords(ctx context.Context) ([]*model.MovieWatchRecord, error) {
	return qs.MovieWatchRecordRepo.FindAll(ctx)
}
