package record

import (
	"context"
	"strconv"
	"strings"
	"time"

	"filmogophery/internal/impression"
	"filmogophery/internal/media"
	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/logger"
)

type (
	CommandService struct {
		MovieWatchRecordRepo ICommandRepository
		WatchMediaRepo       media.IQueryRepository
		ImpressionRepo       impression.ICommandRepository
	}
)

func NewCommandService(
	recordRepo ICommandRepository,
	watchMediaRepo media.IQueryRepository,
	impressionRepo impression.ICommandRepository,
) *CommandService {
	return &CommandService{
		MovieWatchRecordRepo: recordRepo,
		WatchMediaRepo:       watchMediaRepo,
		ImpressionRepo:       impressionRepo,
	}
}

func (cs *CommandService) CreateRecord(dto *CreateMovieRecordDto) error {
	logger := logger.GetLogger()

	code := dto.Media
	watchMediaID, err := cs.WatchMediaRepo.GetMediaIdByCode(context.Background(), &code)
	if err != nil {
		return err
	}

	ctx := context.Background()

	impression := model.MovieImpression{
		ID:      dto.ImpressionID,
		MovieID: dto.MovieID,
		Rating:  &dto.Rating,
		Note:    &dto.Note,
	}
	result, err := cs.ImpressionRepo.UpdateImpression(ctx, &impression)
	if err != nil {
		return err
	}
	logger.Info().Msgf("updated impression: %d", result.RowsAffected)

	date := strings.Split(dto.WatchDate, "-")
	year, _ := strconv.ParseInt(date[0], 10, 64)
	month, _ := strconv.ParseInt(date[1], 10, 64)
	day, _ := strconv.ParseInt(date[2], 10, 64)
	watchRecord := model.MovieWatchRecord{
		WatchMediaID: *watchMediaID,
		WatchDate:    time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.Local),
	}

	_, err = cs.MovieWatchRecordRepo.Save(ctx, &watchRecord)
	if err != nil {
		return err
	}

	return nil
}
