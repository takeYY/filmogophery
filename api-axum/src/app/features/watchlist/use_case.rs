use chrono::NaiveDate;
use serde::{Deserialize, Serialize};

use crate::app::responses::{ApiResult, AppError};

use super::repository::{MovieExistsRepository, WatchlistRepository};

// ─── レスポンス型 ─────────────────────────────────────────────

#[derive(Debug, Serialize)]
pub struct GenreItem {
    pub code: String,
    pub name: String,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct MovieItem {
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
pub struct WatchlistItemResponse {
    pub id: i32,
    pub added_at: Option<chrono::DateTime<chrono::Utc>>,
    pub priority: Option<i32>,
    pub movie: MovieItem,
}

// ─── リクエスト型 ─────────────────────────────────────────────

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct AddToWatchlistInput {
    pub movie_id: i32,
    #[serde(default = "default_priority")]
    pub priority: i32,
}

fn default_priority() -> i32 {
    1
}

// ─── ヘルパー ─────────────────────────────────────────────────

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

/// GET /v1/watchlist — ウォッチリスト一覧取得
pub async fn get_watchlist<WR>(
    repo: &WR,
    user_id: i32,
    limit: i32,
    offset: i32,
) -> ApiResult<Vec<WatchlistItemResponse>>
where
    WR: WatchlistRepository,
{
    let rows = repo.find_by_user_id(user_id, limit, offset).await?;

    let items = rows
        .into_iter()
        .map(|r| WatchlistItemResponse {
            id: r.id,
            added_at: r.added_at,
            priority: r.priority,
            movie: MovieItem {
                id: r.movie_id,
                title: r.movie_title,
                overview: r.movie_overview,
                release_date: r.movie_release_date,
                runtime_minutes: r.movie_runtime_minutes,
                poster_url: r.movie_poster_url,
                tmdb_id: r.movie_tmdb_id,
                genres: parse_genres(r.genre_codes.as_deref(), r.genre_names.as_deref()),
            },
        })
        .collect();

    Ok(items)
}

/// POST /v1/watchlist — ウォッチリスト登録
pub async fn add_to_watchlist<WR, MR>(
    watchlist_repo: &WR,
    movie_repo: &MR,
    user_id: i32,
    input: AddToWatchlistInput,
) -> ApiResult<()>
where
    WR: WatchlistRepository,
    MR: MovieExistsRepository,
{
    // movieId は正の整数であること
    if input.movie_id < 1 {
        let mut errors = std::collections::HashMap::new();
        errors.insert(
            "movieId".to_string(),
            vec!["movieId must be a positive integer".to_string()],
        );
        return Err(AppError::ValidationError(errors));
    }

    // priority は 1〜5 の範囲（Hono のバリデーションに合わせる）
    if input.priority < 1 || input.priority > 5 {
        let mut errors = std::collections::HashMap::new();
        errors.insert(
            "priority".to_string(),
            vec!["priority must be between 1 and 5".to_string()],
        );
        return Err(AppError::ValidationError(errors));
    }

    // 映画の存在確認
    if !movie_repo.exists(input.movie_id).await? {
        return Err(AppError::NotFound("movie".to_string()));
    }

    watchlist_repo
        .create(user_id, input.movie_id, input.priority)
        .await?;

    Ok(())
}

/// DELETE /v1/watchlist/{watchlistId} — ウォッチリストから削除
pub async fn remove_from_watchlist<WR>(
    repo: &WR,
    watchlist_id: i32,
) -> ApiResult<()>
where
    WR: WatchlistRepository,
{
    let affected = repo.delete_by_id(watchlist_id).await?;

    if affected == 0 {
        return Err(AppError::NotFound("watchlist".to_string()));
    }

    Ok(())
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use std::sync::Mutex;

    use chrono::Utc;

    use super::*;
    use crate::app::features::watchlist::repository::{
        MovieExistsRepository, WatchlistRepository, WatchlistRow,
    };
    use crate::app::responses::AppError;

    // ── モック ────────────────────────────────────────────────────

    struct MockWatchlistRepository {
        rows: Vec<WatchlistRow>,
        delete_affected: u64,
        should_fail: bool,
        create_calls: Mutex<u32>,
    }

    impl MockWatchlistRepository {
        fn with_rows(rows: Vec<WatchlistRow>) -> Self {
            Self { rows, delete_affected: 1, should_fail: false, create_calls: Mutex::new(0) }
        }

        fn empty() -> Self {
            Self { rows: vec![], delete_affected: 0, should_fail: false, create_calls: Mutex::new(0) }
        }

        fn failing() -> Self {
            Self { rows: vec![], delete_affected: 0, should_fail: true, create_calls: Mutex::new(0) }
        }

        fn create_count(&self) -> u32 {
            *self.create_calls.lock().unwrap()
        }
    }

    fn sample_row(id: i32, movie_id: i32) -> WatchlistRow {
        WatchlistRow {
            id,
            added_at: Some(Utc::now()),
            priority: Some(1),
            movie_id,
            movie_title: "Test Movie".to_string(),
            movie_overview: "overview".to_string(),
            movie_release_date: None,
            movie_runtime_minutes: 120,
            movie_poster_url: None,
            movie_tmdb_id: movie_id * 10,
            genre_codes: Some("ACTION".to_string()),
            genre_names: Some("アクション".to_string()),
        }
    }

    impl WatchlistRepository for MockWatchlistRepository {
        async fn find_by_user_id(&self, _u: i32, _l: i32, _o: i32) -> Result<Vec<WatchlistRow>, AppError> {
            if self.should_fail { return Err(AppError::InternalServerError); }
            Ok(self.rows.iter().map(|r| WatchlistRow {
                id: r.id,
                added_at: r.added_at,
                priority: r.priority,
                movie_id: r.movie_id,
                movie_title: r.movie_title.clone(),
                movie_overview: r.movie_overview.clone(),
                movie_release_date: r.movie_release_date,
                movie_runtime_minutes: r.movie_runtime_minutes,
                movie_poster_url: r.movie_poster_url.clone(),
                movie_tmdb_id: r.movie_tmdb_id,
                genre_codes: r.genre_codes.clone(),
                genre_names: r.genre_names.clone(),
            }).collect())
        }

        async fn create(&self, _u: i32, _m: i32, _p: i32) -> Result<(), AppError> {
            if self.should_fail { return Err(AppError::InternalServerError); }
            *self.create_calls.lock().unwrap() += 1;
            Ok(())
        }

        async fn delete_by_id(&self, _id: i32) -> Result<u64, AppError> {
            if self.should_fail { return Err(AppError::InternalServerError); }
            Ok(self.delete_affected)
        }
    }

    struct MockMovieRepo { exists: bool }
    impl MockMovieRepo {
        fn found() -> Self { Self { exists: true } }
        fn not_found() -> Self { Self { exists: false } }
    }
    impl MovieExistsRepository for MockMovieRepo {
        async fn exists(&self, _id: i32) -> Result<bool, AppError> { Ok(self.exists) }
    }

    // ── get_watchlist テスト ──────────────────────────────────────

    #[tokio::test]
    async fn test_get_watchlist_returns_mapped_response() {
        let repo = MockWatchlistRepository::with_rows(vec![
            sample_row(1, 10),
            sample_row(2, 20),
        ]);

        let result = get_watchlist(&repo, 1, 12, 0).await.unwrap();

        assert_eq!(result.len(), 2);
        assert_eq!(result[0].id, 1);
        assert_eq!(result[0].movie.id, 10);
        assert_eq!(result[0].movie.genres.len(), 1);
        assert_eq!(result[0].movie.genres[0].code, "ACTION");
    }

    #[tokio::test]
    async fn test_get_watchlist_returns_empty_when_no_data() {
        let repo = MockWatchlistRepository::empty();
        let result = get_watchlist(&repo, 1, 12, 0).await.unwrap();
        assert!(result.is_empty());
    }

    #[tokio::test]
    async fn test_get_watchlist_propagates_repo_error() {
        let repo = MockWatchlistRepository::failing();
        let result = get_watchlist(&repo, 1, 12, 0).await;
        assert!(matches!(result, Err(AppError::InternalServerError)));
    }

    #[tokio::test]
    async fn test_get_watchlist_empty_genres_parsed_correctly() {
        let mut row = sample_row(1, 10);
        row.genre_codes = None;
        row.genre_names = None;
        let repo = MockWatchlistRepository::with_rows(vec![row]);

        let result = get_watchlist(&repo, 1, 12, 0).await.unwrap();

        assert!(result[0].movie.genres.is_empty());
    }

    // ── add_to_watchlist テスト ───────────────────────────────────

    #[tokio::test]
    async fn test_add_to_watchlist_success() {
        let watchlist_repo = MockWatchlistRepository::empty();
        let movie_repo = MockMovieRepo::found();

        let input = AddToWatchlistInput { movie_id: 1, priority: 3 };
        let result = add_to_watchlist(&watchlist_repo, &movie_repo, 1, input).await;

        assert!(result.is_ok());
        assert_eq!(watchlist_repo.create_count(), 1);
    }

    #[tokio::test]
    async fn test_add_to_watchlist_uses_default_priority_1() {
        // default_priority() が 1 を返すことを確認
        let input: AddToWatchlistInput = serde_json::from_str(r#"{"movieId": 1}"#).unwrap();
        assert_eq!(input.priority, 1);
    }

    #[tokio::test]
    async fn test_add_to_watchlist_fails_when_movie_id_invalid() {
        let watchlist_repo = MockWatchlistRepository::empty();
        let movie_repo = MockMovieRepo::found();

        let input = AddToWatchlistInput { movie_id: 0, priority: 1 };
        let result = add_to_watchlist(&watchlist_repo, &movie_repo, 1, input).await;

        assert!(matches!(result, Err(AppError::ValidationError(_))));
    }

    #[tokio::test]
    async fn test_add_to_watchlist_fails_when_priority_out_of_range() {
        let watchlist_repo = MockWatchlistRepository::empty();
        let movie_repo = MockMovieRepo::found();

        let input = AddToWatchlistInput { movie_id: 1, priority: 6 };
        let result = add_to_watchlist(&watchlist_repo, &movie_repo, 1, input).await;

        assert!(matches!(result, Err(AppError::ValidationError(_))));
    }

    #[tokio::test]
    async fn test_add_to_watchlist_fails_when_movie_not_found() {
        let watchlist_repo = MockWatchlistRepository::empty();
        let movie_repo = MockMovieRepo::not_found();

        let input = AddToWatchlistInput { movie_id: 999, priority: 1 };
        let result = add_to_watchlist(&watchlist_repo, &movie_repo, 1, input).await;

        assert!(matches!(result, Err(AppError::NotFound(_))));
        assert_eq!(watchlist_repo.create_count(), 0);
    }

    #[tokio::test]
    async fn test_add_to_watchlist_propagates_repo_error() {
        let watchlist_repo = MockWatchlistRepository::failing();
        let movie_repo = MockMovieRepo::found();

        let input = AddToWatchlistInput { movie_id: 1, priority: 1 };
        let result = add_to_watchlist(&watchlist_repo, &movie_repo, 1, input).await;

        assert!(matches!(result, Err(AppError::InternalServerError)));
    }

    // ── remove_from_watchlist テスト ──────────────────────────────

    #[tokio::test]
    async fn test_remove_from_watchlist_success() {
        let repo = MockWatchlistRepository::with_rows(vec![]);
        // delete_affected = 1 (デフォルト)
        let result = remove_from_watchlist(&repo, 1).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_remove_from_watchlist_not_found() {
        let repo = MockWatchlistRepository::empty();
        // delete_affected = 0
        let result = remove_from_watchlist(&repo, 999).await;
        assert!(matches!(result, Err(AppError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_remove_from_watchlist_propagates_repo_error() {
        let repo = MockWatchlistRepository::failing();
        let result = remove_from_watchlist(&repo, 1).await;
        assert!(matches!(result, Err(AppError::InternalServerError)));
    }

    // ── parse_genres テスト ────────────────────────────────────────

    #[test]
    fn test_parse_genres_multiple() {
        let genres = parse_genres(Some("ACTION,DRAMA"), Some("アクション,ドラマ"));
        assert_eq!(genres.len(), 2);
        assert_eq!(genres[1].name, "ドラマ");
    }

    #[test]
    fn test_parse_genres_none_returns_empty() {
        assert!(parse_genres(None, None).is_empty());
    }
}
