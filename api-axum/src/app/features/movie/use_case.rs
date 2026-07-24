use chrono::NaiveDate;
use serde::{Deserialize, Serialize};

use crate::app::responses::{ApiResult, AppError};
use crate::pkg::tmdb::TmdbClient;

use super::repository::{MovieRepository, NewMovieInput};

// ─── レスポンス型 ─────────────────────────────────────────────

#[derive(Debug, Serialize, Deserialize)]
pub struct GenreItem {
    pub code: String,
    pub name: String,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct MovieResponse {
    pub id: i32,
    pub title: String,
    pub overview: String,
    pub release_date: Option<NaiveDate>,
    pub runtime_minutes: i32,
    pub poster_url: Option<String>,
    pub tmdb_id: i32,
    pub genres: Vec<GenreItem>,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct ReviewItem {
    pub id: i32,
    pub rating: f64,
    pub comment: Option<String>,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub updated_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct SeriesItem {
    pub name: String,
    pub poster_url: Option<String>,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct MovieDetailResponse {
    pub id: i32,
    pub title: String,
    pub overview: String,
    pub release_date: Option<NaiveDate>,
    pub runtime_minutes: i32,
    pub poster_url: Option<String>,
    pub tmdb_id: i32,
    pub genres: Vec<GenreItem>,
    pub vote_average: f64,
    pub vote_count: i32,
    pub series: Option<SeriesItem>,
    pub review: Option<ReviewItem>,
}

// ─── ヘルパー ─────────────────────────────────────────────────

/// "CODE1,CODE2" と "NAME1,NAME2" からジャンルリストを構築する
fn parse_genres(codes: Option<&str>, names: Option<&str>) -> Vec<GenreItem> {
    let codes: Vec<&str> = codes
        .filter(|s| !s.is_empty())
        .map(|s| s.split(',').collect())
        .unwrap_or_default();
    let names: Vec<&str> = names
        .filter(|s| !s.is_empty())
        .map(|s| s.split(',').collect())
        .unwrap_or_default();

    codes
        .into_iter()
        .zip(names.into_iter())
        .map(|(code, name)| GenreItem {
            code: code.to_string(),
            name: name.to_string(),
        })
        .collect()
}

// ─── Use Case 関数 ─────────────────────────────────────────────

/// GET /v1/movies — ユーザーがレビューした映画一覧
pub async fn get_movies<R>(
    repo: &R,
    user_id: i32,
    genre: Option<&str>,
    limit: i32,
    offset: i32,
) -> ApiResult<Vec<MovieResponse>>
where
    R: MovieRepository,
{
    let rows = repo.find_reviewed_by_user(user_id, genre, limit, offset).await?;

    let movies = rows
        .into_iter()
        .map(|r| MovieResponse {
            id: r.id,
            title: r.title,
            overview: r.overview,
            release_date: r.release_date,
            runtime_minutes: r.runtime_minutes,
            poster_url: r.poster_url,
            tmdb_id: r.tmdb_id,
            genres: parse_genres(r.genre_codes.as_deref(), r.genre_names.as_deref()),
        })
        .collect();

    Ok(movies)
}

/// GET /v1/movies/:id — 映画詳細（TMDB API 連携）
pub async fn get_movie_detail<R>(
    repo: &R,
    tmdb: &TmdbClient,
    movie_id: i32,
    user_id: i32,
) -> ApiResult<MovieDetailResponse>
where
    R: MovieRepository,
{
    let row = repo
        .find_detail_by_id(movie_id, user_id)
        .await?
        .ok_or_else(|| AppError::NotFound("movie".to_string()))?;

    // TMDB から vote_average / vote_count / runtime を取得
    let tmdb_detail = tmdb.get_movie_detail(row.tmdb_id).await?;

    // runtime_minutes が 0 の場合は TMDB の値で更新
    let runtime = if row.runtime_minutes == 0 {
        if let Some(rt) = tmdb_detail.runtime {
            let _ = repo.update_runtime_minutes(movie_id, rt).await;
            rt
        } else {
            row.runtime_minutes
        }
    } else {
        row.runtime_minutes
    };

    let genres = parse_genres(row.genre_codes.as_deref(), row.genre_names.as_deref());

    let review = row.review_id.map(|id| ReviewItem {
        id,
        rating: row.review_rating.unwrap_or(0.0),
        comment: row.review_comment,
        created_at: row.review_created_at.unwrap_or_default(),
        updated_at: row.review_updated_at.unwrap_or_default(),
    });

    let series = row.series_name.map(|name| SeriesItem {
        name,
        poster_url: row.series_poster_url,
    });

    Ok(MovieDetailResponse {
        id: row.id,
        title: row.title,
        overview: row.overview,
        release_date: row.release_date,
        runtime_minutes: runtime,
        poster_url: row.poster_url,
        tmdb_id: row.tmdb_id,
        genres,
        vote_average: tmdb_detail.vote_average_5(),
        vote_count: tmdb_detail.vote_count,
        series,
        review,
    })
}

/// GET /v1/search/movies — TMDB 映画検索（新規映画は DB に保存、結果を Redis に 24 時間キャッシュ）
pub async fn search_movies<R>(
    repo: &R,
    tmdb: &TmdbClient,
    redis: Option<&crate::pkg::redis::RedisClient>,
    title: &str,
    limit: i32,
    offset: i32,
) -> ApiResult<Vec<MovieResponse>>
where
    R: MovieRepository,
{
    // ─── Redis キャッシュチェック ──────────────────────────────
    let cache_key = format!(
        "movies:search:{}:limit:{}:offset:{}",
        title.trim().to_lowercase(),
        limit,
        offset
    );

    if let Some(redis) = redis {
        if let Some(cached) = redis.get::<Vec<MovieResponse>>(&cache_key).await {
            tracing::debug!("cache hit: {cache_key}");
            return Ok(cached);
        }
    }

    // ─── TMDB API 検索 ────────────────────────────────────────
    let tmdb_result = tmdb.search_movies(title, 1).await?;
    let tmdb_movies = tmdb_result.results;

    if tmdb_movies.is_empty() {
        return Ok(vec![]);
    }

    // TMDB ID リストで既存映画を取得
    let tmdb_ids: Vec<i32> = tmdb_movies.iter().map(|m| m.id).collect();
    let existing = repo.find_by_tmdb_ids(&tmdb_ids).await?;
    let existing_tmdb_ids: std::collections::HashSet<i32> =
        existing.iter().map(|r| r.tmdb_id).collect();

    // DB に存在しない映画を一括登録
    let new_movies: Vec<NewMovieInput> = tmdb_movies
        .iter()
        .filter(|m| !existing_tmdb_ids.contains(&m.id))
        .map(|m| NewMovieInput {
            tmdb_id: m.id,
            title: m.title.clone(),
            overview: m.overview.clone(),
            release_date: if m.release_date.is_empty() {
                "1970-01-01".to_string()
            } else {
                m.release_date.clone()
            },
            poster_url: m.poster_path.clone(),
            genre_ids: m.genre_ids.clone(),
        })
        .collect();

    if !new_movies.is_empty() {
        repo.batch_insert(&new_movies).await?;
    }

    // 新規登録分を含めて再取得
    let all_rows = repo.find_by_tmdb_ids(&tmdb_ids).await?;
    let row_map: std::collections::HashMap<i32, &super::repository::MovieTmdbRow> =
        all_rows.iter().map(|r| (r.tmdb_id, r)).collect();

    // TMDB の順序を保ちながら offset/limit を適用してレスポンス構築
    let results: Vec<MovieResponse> = tmdb_movies
        .iter()
        .skip(offset as usize)
        .take(limit as usize)
        .filter_map(|t| {
            let row = row_map.get(&t.id)?;
            Some(MovieResponse {
                id: row.id,
                title: t.title.clone(),
                overview: t.overview.clone(),
                release_date: NaiveDate::parse_from_str(&t.release_date, "%Y-%m-%d").ok(),
                runtime_minutes: 0,
                poster_url: t.poster_path.clone(),
                tmdb_id: t.id,
                genres: parse_genres(row.genre_codes.as_deref(), row.genre_names.as_deref()),
            })
        })
        .collect();

    // ─── Redis にキャッシュ（24 時間）────────────────────────
    if let Some(redis) = redis {
        redis.set(&cache_key, &results, 24 * 60 * 60).await;
    }

    Ok(results)
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use chrono::Utc;
    use std::sync::Mutex;

    use super::*;
    use crate::app::features::movie::repository::{
        MovieDetailRow, MovieRepository, MovieRow, MovieTmdbRow, NewMovieInput,
    };
    use crate::app::responses::AppError;

    // ── モック: MovieRepository ───────────────────────────────────

    struct MockMovieRepository {
        reviewed_rows: Vec<MovieRow>,
        detail_row: Option<MovieDetailRow>,
        tmdb_rows: Vec<MovieTmdbRow>,
        batch_insert_count: Mutex<u32>,
        should_fail: bool,
    }

    impl MockMovieRepository {
        fn with_reviewed(rows: Vec<MovieRow>) -> Self {
            Self {
                reviewed_rows: rows,
                detail_row: None,
                tmdb_rows: vec![],
                batch_insert_count: Mutex::new(0),
                should_fail: false,
            }
        }

        fn with_detail(row: MovieDetailRow) -> Self {
            Self {
                reviewed_rows: vec![],
                detail_row: Some(row),
                tmdb_rows: vec![],
                batch_insert_count: Mutex::new(0),
                should_fail: false,
            }
        }

        fn empty() -> Self {
            Self {
                reviewed_rows: vec![],
                detail_row: None,
                tmdb_rows: vec![],
                batch_insert_count: Mutex::new(0),
                should_fail: false,
            }
        }

        fn failing() -> Self {
            Self {
                reviewed_rows: vec![],
                detail_row: None,
                tmdb_rows: vec![],
                batch_insert_count: Mutex::new(0),
                should_fail: true,
            }
        }

        fn insert_count(&self) -> u32 {
            *self.batch_insert_count.lock().unwrap()
        }
    }

    fn sample_movie_row(id: i32, title: &str) -> MovieRow {
        MovieRow {
            id,
            title: title.to_string(),
            overview: "overview".to_string(),
            release_date: None,
            runtime_minutes: 120,
            poster_url: None,
            tmdb_id: id * 100,
            genre_codes: Some("ACTION".to_string()),
            genre_names: Some("アクション".to_string()),
        }
    }

    fn sample_detail_row(id: i32, tmdb_id: i32, with_review: bool) -> MovieDetailRow {
        MovieDetailRow {
            id,
            title: "Test Movie".to_string(),
            overview: "overview".to_string(),
            release_date: None,
            runtime_minutes: 0, // TMDB から取得するケースをテスト
            poster_url: None,
            tmdb_id,
            genre_codes: Some("DRAMA".to_string()),
            genre_names: Some("ドラマ".to_string()),
            review_id: if with_review { Some(10) } else { None },
            review_rating: if with_review { Some(4.0) } else { None },
            review_comment: if with_review { Some("good".to_string()) } else { None },
            review_created_at: if with_review { Some(Utc::now()) } else { None },
            review_updated_at: if with_review { Some(Utc::now()) } else { None },
            series_name: None,
            series_poster_url: None,
        }
    }

    impl MovieRepository for MockMovieRepository {
        async fn find_reviewed_by_user(
            &self,
            _user_id: i32,
            _genre: Option<&str>,
            _limit: i32,
            _offset: i32,
        ) -> Result<Vec<MovieRow>, AppError> {
            if self.should_fail {
                return Err(AppError::InternalServerError);
            }
            Ok(self.reviewed_rows.iter().map(|r| MovieRow {
                id: r.id,
                title: r.title.clone(),
                overview: r.overview.clone(),
                release_date: r.release_date,
                runtime_minutes: r.runtime_minutes,
                poster_url: r.poster_url.clone(),
                tmdb_id: r.tmdb_id,
                genre_codes: r.genre_codes.clone(),
                genre_names: r.genre_names.clone(),
            }).collect())
        }

        async fn find_detail_by_id(
            &self,
            _movie_id: i32,
            _user_id: i32,
        ) -> Result<Option<MovieDetailRow>, AppError> {
            if self.should_fail {
                return Err(AppError::InternalServerError);
            }
            Ok(self.detail_row.as_ref().map(|r| MovieDetailRow {
                id: r.id,
                title: r.title.clone(),
                overview: r.overview.clone(),
                release_date: r.release_date,
                runtime_minutes: r.runtime_minutes,
                poster_url: r.poster_url.clone(),
                tmdb_id: r.tmdb_id,
                genre_codes: r.genre_codes.clone(),
                genre_names: r.genre_names.clone(),
                review_id: r.review_id,
                review_rating: r.review_rating,
                review_comment: r.review_comment.clone(),
                review_created_at: r.review_created_at,
                review_updated_at: r.review_updated_at,
                series_name: r.series_name.clone(),
                series_poster_url: r.series_poster_url.clone(),
            }))
        }

        async fn find_by_tmdb_ids(&self, _tmdb_ids: &[i32]) -> Result<Vec<MovieTmdbRow>, AppError> {
            if self.should_fail {
                return Err(AppError::InternalServerError);
            }
            Ok(self.tmdb_rows.iter().map(|r| MovieTmdbRow {
                id: r.id,
                tmdb_id: r.tmdb_id,
                genre_codes: r.genre_codes.clone(),
                genre_names: r.genre_names.clone(),
            }).collect())
        }

        async fn update_runtime_minutes(&self, _movie_id: i32, _runtime: i32) -> Result<(), AppError> {
            Ok(())
        }

        async fn batch_insert(&self, _movies: &[NewMovieInput]) -> Result<(), AppError> {
            if self.should_fail {
                return Err(AppError::InternalServerError);
            }
            *self.batch_insert_count.lock().unwrap() += 1;
            Ok(())
        }
    }

    // ── get_movies テスト ─────────────────────────────────────────

    #[tokio::test]
    async fn test_get_movies_returns_mapped_response() {
        let repo = MockMovieRepository::with_reviewed(vec![
            sample_movie_row(1, "映画A"),
            sample_movie_row(2, "映画B"),
        ]);

        let result = get_movies(&repo, 1, None, 12, 0).await.unwrap();

        assert_eq!(result.len(), 2);
        assert_eq!(result[0].id, 1);
        assert_eq!(result[0].title, "映画A");
        assert_eq!(result[0].genres.len(), 1);
        assert_eq!(result[0].genres[0].code, "ACTION");
    }

    #[tokio::test]
    async fn test_get_movies_returns_empty_when_no_reviews() {
        let repo = MockMovieRepository::empty();
        let result = get_movies(&repo, 1, None, 12, 0).await.unwrap();
        assert!(result.is_empty());
    }

    #[tokio::test]
    async fn test_get_movies_propagates_repo_error() {
        let repo = MockMovieRepository::failing();
        let result = get_movies(&repo, 1, None, 12, 0).await;
        assert!(matches!(result, Err(AppError::InternalServerError)));
    }

    #[tokio::test]
    async fn test_get_movies_genres_parsed_correctly() {
        let mut row = sample_movie_row(1, "映画");
        row.genre_codes = Some("ACTION,DRAMA".to_string());
        row.genre_names = Some("アクション,ドラマ".to_string());
        let repo = MockMovieRepository::with_reviewed(vec![row]);

        let result = get_movies(&repo, 1, None, 12, 0).await.unwrap();

        assert_eq!(result[0].genres.len(), 2);
        assert_eq!(result[0].genres[1].code, "DRAMA");
    }

    #[tokio::test]
    async fn test_get_movies_empty_genre_codes_returns_empty_genres() {
        let mut row = sample_movie_row(1, "映画");
        row.genre_codes = None;
        row.genre_names = None;
        let repo = MockMovieRepository::with_reviewed(vec![row]);

        let result = get_movies(&repo, 1, None, 12, 0).await.unwrap();

        assert!(result[0].genres.is_empty());
    }

    // ── get_movie_detail テスト (DB 部分のみ: TMDB は別途結合テスト) ──

    #[tokio::test]
    async fn test_get_movie_detail_returns_not_found_when_missing() {
        let repo = MockMovieRepository::empty();
        // TmdbClient をダミー（テストではリクエストが走らないケースのみ確認）
        // → NotFound は DB 取得後に早期リターンするため TMDB 呼び出しは起きない
        let tmdb = TmdbClient::new("dummy");

        let result = get_movie_detail(&repo, &tmdb, 999, 1).await;

        assert!(matches!(result, Err(AppError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_get_movie_detail_propagates_db_error() {
        let repo = MockMovieRepository::failing();
        let tmdb = TmdbClient::new("dummy");

        let result = get_movie_detail(&repo, &tmdb, 1, 1).await;

        assert!(matches!(result, Err(AppError::InternalServerError)));
    }

    // ── parse_genres テスト ────────────────────────────────────────

    #[test]
    fn test_parse_genres_single() {
        let genres = parse_genres(Some("ACTION"), Some("アクション"));
        assert_eq!(genres.len(), 1);
        assert_eq!(genres[0].code, "ACTION");
        assert_eq!(genres[0].name, "アクション");
    }

    #[test]
    fn test_parse_genres_multiple() {
        let genres = parse_genres(Some("ACTION,DRAMA"), Some("アクション,ドラマ"));
        assert_eq!(genres.len(), 2);
    }

    #[test]
    fn test_parse_genres_none_returns_empty() {
        assert!(parse_genres(None, None).is_empty());
    }

    #[test]
    fn test_parse_genres_empty_string_returns_empty() {
        assert!(parse_genres(Some(""), Some("")).is_empty());
    }
}
