package movie

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/config"
	"filmogophery/internal/db"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/tmdb"
)

type (
	CommandService struct {
		MovieRepo      repositories.IMovieRepository
		GenreRepo      repositories.IGenreRepository
		ImpressionRepo repositories.IImpressionRepository
		TmdbClient     tmdb.ITmdbClient
	}
)

func NewCommandService(
	movieRepo repositories.IMovieRepository,
	genreRepo repositories.IGenreRepository,
	impressionRepo repositories.IImpressionRepository,
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
	conf := config.LoadConfig()
	gormDB := db.ConnectDB(conf)

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

	genres, err := cs.GenreRepo.FindByNames(context.Background(), genreNames)
	if err != nil {
		return err
	}

	movie := model.Movie{
		ID:          dto.TmdbID,
		Title:       tmdbMovie.Title,
		Overview:    tmdbMovie.Overview,
		ReleaseDate: time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, time.Local),
		RunTime:     int32(tmdbMovie.Runtime),
		PosterURL:   &tmdbMovie.PosterPath,
		TmdbID:      dto.TmdbID,
		Genres:      genres,
	}

	impression := model.MovieImpression{
		MovieID: dto.TmdbID,
		Status:  dto.Status,
	}

	ctx := context.Background()
	err = gormDB.Transaction(func(tx *gorm.DB) error {
		if e := cs.MovieRepo.Save(ctx, repositories.SaveMovieInput{Target: &movie}); e != nil {
			return e
		}
		if e := cs.ImpressionRepo.Save(ctx, repositories.SaveImpressionInput{Target: &impression}); e != nil {
			return e
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
