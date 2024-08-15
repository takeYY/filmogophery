package record

import (
	"context"

	"filmogophery/internal/db"
	"filmogophery/pkg/gen/model"
)

type (
	QueryService struct {
		MovieWatchRecordRepo IQueryRepository
	}
)

func NewQueryService() *QueryService {
	var movieWatchRecordRepo IQueryRepository = &MovieWatchRecordRepository{
		DB: db.READER_DB,
	}

	return &QueryService{
		MovieWatchRecordRepo: movieWatchRecordRepo,
	}
}

func (qs *QueryService) GetWatchRecords(ctx context.Context) ([]*model.MovieWatchRecord, error) {
	return qs.MovieWatchRecordRepo.Find(ctx)
}
