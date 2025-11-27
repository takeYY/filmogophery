package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/config"
)

type (
	ITmdbRepository interface {
		// IDに一致する映画詳細を取得
		GetMovieDetail(id int32) (*types.TmdbMovieDetail, error)
		// タイトルに一致する映画一覧を取得
		GetMoviesByTitle(title string, page int32) (*types.TmdbSearchMovieResult, error)
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

// タイトルに一致する映画一覧を取得
func (r *tmdbRepository) GetMoviesByTitle(title string, page int32) (*types.TmdbSearchMovieResult, error) {
	URL := r.baseURL + "/search/movie"
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", title)
	q.Set("page", strconv.Itoa(int(page)))
	q.Set("language", "ja-JP")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

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

	var result types.TmdbSearchMovieResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
