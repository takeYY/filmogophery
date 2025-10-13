package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"filmogophery/internal/app/types"
	"filmogophery/internal/config"
)

type (
	ITmdbRepository interface {
		// IDに一致する映画詳細を取得
		GetMovieDetail(id int32) (*types.TmdbMovieDetail, error)
	}

	tmdbRepository struct {
		httpClient  *http.Client
		baseURL     string
		AccessToken string
	}
)

func NewTmdbRepository(tmdb *config.Tmdb) ITmdbRepository {
	return &tmdbRepository{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL:     "https://api.themoviedb.org/3",
		AccessToken: tmdb.ACCESS_TOKEN,
	}
}

// IDに一致する映画詳細を取得
func (r *tmdbRepository) GetMovieDetail(id int32) (*types.TmdbMovieDetail, error) {
	url := fmt.Sprintf("%s/movie/%d?language=ja-JP", r.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+r.AccessToken)

	res, err := r.httpClient.Do(req)
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

	var result types.TmdbMovieDetail
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
