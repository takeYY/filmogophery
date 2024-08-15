package movie

import (
	"context"
	"strconv"
	"strings"
	"time"

	"filmogophery/internal/db"
	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/logger"
)

type (
	CommandService struct {
		MovieRepo ICommandRepository
	}
)

func NewCommandService() *CommandService {
	var movieRepo ICommandRepository = &MovieRepository{
		DB: db.WRITER_DB,
	}

	return &CommandService{
		MovieRepo: movieRepo,
	}
}

func (cs *CommandService) CreateMovie(dto *CreateMovieDto) (*model.Movie, error) {
	date := strings.Split(dto.ReleaseDate, "-")
	year, _ := strconv.ParseInt(date[0], 10, 64)
	month, _ := strconv.ParseInt(date[1], 10, 64)
	day, _ := strconv.ParseInt(date[2], 10, 64)

	movie := model.Movie{
		Title:       dto.Title,
		Overview:    dto.Overview,
		ReleaseDate: time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.Local),
		RunTime:     dto.RunTime,
		TmdbID:      dto.TmdbID,
		Genres:      []*model.Genre{},
	}

	return cs.MovieRepo.Save(&movie)
}

func (cs *CommandService) CreateMovieRecord(dto *CreateMovieRecordDto) error {
	logger := logger.GetLogger()

	code := dto.Media
	watchMediaID, err := cs.MovieRepo.GetMediaIdByCode(context.Background(), &code)
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
	result, err := cs.MovieRepo.UpdateImpression(ctx, &impression)
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

	_, e := cs.MovieRepo.SaveRecord(ctx, &watchRecord)
	if e != nil {
		return e
	}

	return nil
}
