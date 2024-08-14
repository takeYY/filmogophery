package tmdb

import (
	"strings"

	"filmogophery/internal/config"
	"filmogophery/pkg/logger"
	"filmogophery/pkg/tokenizer"
)

type (
	TmdbService struct {
		tmdbClient ITmdbClient
	}
)

func NewTmdbService(conf *config.Config) *TmdbService {
	var tmdbClient ITmdbClient = NewTmdbClient(conf)
	return &TmdbService{
		tmdbClient: tmdbClient,
	}
}

func (ts *TmdbService) SearchMovies(q *string) (*SearchMovieResultSet, error) {
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

	return movies, nil
}
