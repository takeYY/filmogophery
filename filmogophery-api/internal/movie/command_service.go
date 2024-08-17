package movie

import (
	"context"
	"strconv"
	"strings"
	"time"

	"filmogophery/internal/db"
	"filmogophery/internal/impression"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	CommandService struct {
		MovieRepo      ICommandRepository
		ImpressionRepo impression.ICommandRepository
	}
)

func NewCommandService(movieRepo ICommandRepository, impressionRepo impression.ICommandRepository) *CommandService {
	return &CommandService{
		MovieRepo:      movieRepo,
		ImpressionRepo: impressionRepo,
	}
}

func (cs *CommandService) CreateMovieAndImpression(dto *CreateMovieDto) error {
	date := strings.Split(dto.ReleaseDate, "-")
	year, _ := strconv.ParseInt(date[0], 10, 64)
	month, _ := strconv.ParseInt(date[1], 10, 64)
	day, _ := strconv.ParseInt(date[2], 10, 64)

	movie := model.Movie{
		ID:          dto.TmdbID,
		Title:       dto.Title,
		Overview:    dto.Overview,
		ReleaseDate: time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.Local),
		RunTime:     dto.RunTime,
		PosterURL:   dto.PosterURL,
		TmdbID:      &dto.TmdbID,
		Genres:      []*model.Genre{},
	}

	impression := model.MovieImpression{
		MovieID: dto.TmdbID,
		Status:  dto.Status,
	}

	q := query.Use(db.WRITER_DB)
	ctx := context.Background()
	q.Transaction(func(tx *query.Query) error {
		if _, err := cs.MovieRepo.Save(ctx, &movie); err != nil {
			return err
		}
		if _, err := cs.ImpressionRepo.Save(ctx, &impression); err != nil {
			return err
		}

		return nil
	})

	return nil
}
