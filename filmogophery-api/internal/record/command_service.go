package record

import (
	"context"
	"strconv"
	"strings"
	"time"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	CommandService struct {
		MovieWatchRecordRepo repositories.IRecordRepository
		WatchMediaRepo       repositories.IMediaRepository
		ImpressionRepo       repositories.IImpressionRepository
	}
)

func NewCommandService(
	recordRepo repositories.IRecordRepository,
	watchMediaRepo repositories.IMediaRepository,
	impressionRepo repositories.IImpressionRepository,
) *CommandService {
	return &CommandService{
		recordRepo,
		watchMediaRepo,
		impressionRepo,
	}
}

func (cs *CommandService) CreateRecord(dto *CreateMovieRecordDto) error {
	logger := logger.GetLogger()

	code := dto.Media
	watchMedia, err := cs.WatchMediaRepo.FindByCode(context.Background(), &code)
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
	affected, err := cs.ImpressionRepo.Update(ctx, repositories.UpdateImpressionInput{Target: &impression})
	if err != nil {
		return err
	}
	logger.Info().Msgf("updated impression: %d", affected)

	date := strings.Split(dto.WatchDate, "-")
	year, _ := strconv.ParseInt(date[0], 10, 64)
	month, _ := strconv.ParseInt(date[1], 10, 64)
	day, _ := strconv.ParseInt(date[2], 10, 64)
	watchRecord := model.MovieWatchRecord{
		WatchMediaID: watchMedia.ID,
		WatchDate:    time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.Local),
	}

	err = cs.MovieWatchRecordRepo.Save(ctx, repositories.SaveRecordInput{Target: &watchRecord})
	if err != nil {
		return err
	}

	return nil
}
