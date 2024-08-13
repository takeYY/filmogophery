package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"filmogophery/internal/config"
)

type (
	ITmdbClient interface {
		SearchMovies(query string) (*SearchMovieResultSet, error)
	}

	TmdbClient struct {
		httpClient  *http.Client
		baseURL     string
		AccessToken string
	}
)

func NewTmdbClient(conf *config.Config) *TmdbClient {
	return &TmdbClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL:     "https://api.themoviedb.org/3",
		AccessToken: conf.Tmdb.ACCESS_TOKEN,
	}
}

func (tc *TmdbClient) SearchMovies(query string) (*SearchMovieResultSet, error) {
	url := fmt.Sprintf("%s/search/movie?query=%s&language=ja-JP", tc.baseURL, url.QueryEscape(query))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tc.AccessToken)

	res, err := tc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch movies: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result SearchMovieResultSet
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
