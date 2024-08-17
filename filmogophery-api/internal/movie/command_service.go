package movie

import (
	"context"
	"strconv"
	"strings"
	"time"

	"filmogophery/internal/db"
	"filmogophery/internal/genre"
	"filmogophery/internal/impression"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
	"filmogophery/internal/tmdb"
)

type (
	CommandService struct {
		MovieRepo      ICommandRepository
		GenreRepo      genre.IQueryRepository
		ImpressionRepo impression.ICommandRepository
		TmdbClient     tmdb.ITmdbClient
	}
)

func NewCommandService(
	movieRepo ICommandRepository,
	genreRepo genre.IQueryRepository,
	impressionRepo impression.ICommandRepository,
	tmdbClient tmdb.ITmdbClient,
) *CommandService {
	return &CommandService{
		MovieRepo:      movieRepo,
		GenreRepo:      genreRepo,
		ImpressionRepo: impressionRepo,
		TmdbClient:     tmdbClient,
	}
}

func (cs *CommandService) CreateMovieAndImpression(dto *CreateMovieDto) error {
	tmdbMovie, err := cs.TmdbClient.GetMovieDetail(&dto.TmdbID)
	if err != nil {
		return err
	}

	date := strings.Split(tmdbMovie.ReleaseDate, "-")
	year, _ := strconv.ParseInt(date[0], 10, 64)
	month, _ := strconv.ParseInt(date[1], 10, 64)
	day, _ := strconv.ParseInt(date[2], 10, 64)

	var genreNames []string = make([]string, 0)
	for _, g := range tmdbMovie.Genres {
		genre := tmdb.GetGenreName(&g.ID)
		genreNames = append(genreNames, genre)
	}

	genres, err := cs.GenreRepo.FindByName(context.Background(), genreNames)
	if err != nil {
		return err
	}

	movie := model.Movie{
		ID:          dto.TmdbID,
		Title:       tmdbMovie.Title,
		Overview:    &tmdbMovie.Overview,
		ReleaseDate: time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.Local),
		RunTime:     int32(tmdbMovie.Runtime),
		PosterURL:   &tmdbMovie.PosterPath,
		TmdbID:      &dto.TmdbID,
		Genres:      genres,
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
