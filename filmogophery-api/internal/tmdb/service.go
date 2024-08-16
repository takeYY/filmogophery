package tmdb

import (
	"strings"

	"filmogophery/pkg/logger"
	"filmogophery/pkg/tokenizer"
)

type (
	TmdbService struct {
		tmdbClient ITmdbClient
	}
)

func NewTmdbService(tmdbClient ITmdbClient) *TmdbService {
	return &TmdbService{
		tmdbClient: tmdbClient,
	}
}

func (ts *TmdbService) SearchMovies(q *string) ([]*SearchMovieDto, error) {
	logger := logger.GetLogger()

	ch := make(chan *tokenizer.NEologd, 1)
	r := strings.NewReader(*q)

	go tokenizer.SyncTokenize(ch, r)

	var qs []string
	for {
		k, ok := <-ch
		if !ok {
			break
		}
		qs = append(qs, k.Surface)
	}
	query := strings.Join(qs, " ")
	logger.Info().Msgf("query is [%s]", query)

	movies, err := ts.tmdbClient.SearchMovies(query)
	if err != nil {
		logger.Error().Msg("Error fetching movies")
		return nil, err
	}

	var results []*SearchMovieDto = make([]*SearchMovieDto, 0)
	for _, movie := range movies.Results {
		var genres []*string
		for _, genre := range movie.GenreIds {
			g := GetGenreName(genre)
			genres = append(genres, &g)
		}
		result := &SearchMovieDto{
			TmdbID:      int32(movie.ID),
			Title:       movie.Title,
			Overview:    *movie.Overview,
			Popularity:  float32(movie.Popularity),
			PosterURL:   *movie.PosterPath,
			ReleaseDate: movie.ReleaseDate,
			VoteAverage: float32(movie.VoteAverage),
			VoteCount:   int32(movie.VoteCount),
			Genres:      genres,
		}
		results = append(results, result)
	}

	return results, nil
}
