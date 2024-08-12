package movie

import (
	"filmogophery/internal/db"
	"filmogophery/pkg/gen/model"
	"strconv"
	"strings"
	"time"
)

type (
	CommandService struct {
		MovieRepo ICommandRepository
	}
)

func NewCommandService() *CommandService {
	var movieRepo ICommandRepository = &MovieRepository{
		DB: db.READER_DB,
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
