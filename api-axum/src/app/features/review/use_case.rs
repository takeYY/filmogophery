use serde::Deserialize;

use crate::app::responses::{ApiResult, AppError};

use super::repository::{
    MovieExistsRepository, PlatformExistsRepository, PointRepository, ReviewRepository,
    WatchHistoryRepository,
};

// ─── 定数 ────────────────────────────────────────────────────

const POINTS_FOR_REVIEW: i32 = 20;

/// 上映時間からポイントを計算する（Echo/Hono と同じロジック）
fn calc_watch_points(runtime_minutes: i32) -> i32 {
    if runtime_minutes <= 90 {
        10
    } else if runtime_minutes <= 150 {
        15
    } else {
        20
    }
}

// ─── リクエスト型 ─────────────────────────────────────────────

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WatchHistoryInput {
    pub platform_id: i32,
    pub watched_date: Option<String>, // "YYYY-MM-DD"
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CreateReviewInput {
    pub rating: Option<f64>,
    pub comment: Option<String>,
    pub watch_history: Option<WatchHistoryInput>,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct UpdateReviewInput {
    pub rating: Option<f64>,
    pub comment: Option<String>,
}

// ─── Use Case 関数 ─────────────────────────────────────────────

/// POST /v1/movies/{movieId}/reviews — レビュー登録
///
/// - 映画の存在確認
/// - レビュー重複チェック
/// - プラットフォームの存在確認（視聴履歴あり時のみ）
/// - レビュー作成 → レビューポイント付与
/// - 視聴履歴作成 → 視聴履歴ポイント付与（視聴履歴あり時のみ）
///
/// NOTE: トランザクション制御は sqlx の begin/commit が必要だが、
///       ポイント付与を含む複合処理のため各操作を逐次実行する。
///       部分失敗時のロールバックは将来の課題（現在 Echo/Hono も同様の構成）。
pub async fn create_review<RR, MR, PR, WR, PT>(
    review_repo: &RR,
    movie_repo: &MR,
    platform_repo: &PR,
    watch_history_repo: &WR,
    point_repo: &PT,
    user_id: i32,
    movie_id: i32,
    input: CreateReviewInput,
) -> ApiResult<()>
where
    RR: ReviewRepository,
    MR: MovieExistsRepository,
    PR: PlatformExistsRepository,
    WR: WatchHistoryRepository,
    PT: PointRepository,
{
    // rating か comment のどちらかは必須
    if input.rating.is_none() && input.comment.is_none() {
        let mut errors = std::collections::HashMap::new();
        errors.insert(
            "rating".to_string(),
            vec!["rating or comment is required".to_string()],
        );
        return Err(AppError::ValidationError(errors));
    }

    // rating のバリデーション
    if let Some(r) = input.rating {
        if r < 0.1 || r > 5.0 {
            let mut errors = std::collections::HashMap::new();
            errors.insert(
                "rating".to_string(),
                vec!["rating must be between 0.1 and 5.0".to_string()],
            );
            return Err(AppError::ValidationError(errors));
        }
    }

    // 映画の存在確認
    let movie = movie_repo
        .find_by_id(movie_id)
        .await?
        .ok_or_else(|| AppError::NotFound("movie".to_string()))?;

    // レビュー重複チェック
    if review_repo.find_by_movie_id(user_id, movie_id).await?.is_some() {
        return Err(AppError::Conflict("review".to_string()));
    }

    // プラットフォームの存在確認（視聴履歴あり時のみ）
    if let Some(wh) = &input.watch_history {
        if !platform_repo.exists(wh.platform_id).await? {
            return Err(AppError::NotFound("platform".to_string()));
        }
    }

    // レビューを作成
    let review_id = review_repo
        .create(
            user_id,
            movie_id,
            input.rating,
            input.comment.as_deref(),
        )
        .await?;

    // レビューポイント付与
    point_repo
        .grant_points(user_id, POINTS_FOR_REVIEW, "review", review_id)
        .await?;

    // 視聴履歴を作成（入力がある場合のみ）
    if let Some(wh_input) = &input.watch_history {
        let wh_id = watch_history_repo
            .create(
                user_id,
                movie_id,
                wh_input.platform_id,
                wh_input.watched_date.as_deref(),
            )
            .await?;

        // 視聴履歴ポイント付与
        let watch_points = calc_watch_points(movie.runtime_minutes);
        point_repo
            .grant_points(user_id, watch_points, "watch_history", wh_id)
            .await?;
    }

    Ok(())
}

/// PUT /v1/reviews/{reviewId} — レビュー更新
///
/// - レビューの存在確認（ユーザー所有チェック込み）
/// - rating か comment の少なくとも一方を更新
pub async fn update_review<RR>(
    review_repo: &RR,
    user_id: i32,
    review_id: i32,
    input: UpdateReviewInput,
) -> ApiResult<()>
where
    RR: ReviewRepository,
{
    // rating か comment のどちらかは必須
    if input.rating.is_none() && input.comment.is_none() {
        let mut errors = std::collections::HashMap::new();
        errors.insert(
            "rating".to_string(),
            vec!["rating or comment is required".to_string()],
        );
        return Err(AppError::ValidationError(errors));
    }

    // rating のバリデーション
    if let Some(r) = input.rating {
        if r < 0.1 || r > 5.0 {
            let mut errors = std::collections::HashMap::new();
            errors.insert(
                "rating".to_string(),
                vec!["rating must be between 0.1 and 5.0".to_string()],
            );
            return Err(AppError::ValidationError(errors));
        }
    }

    // レビューの存在確認（ユーザー所有チェック込み）
    review_repo
        .find_by_id(user_id, review_id)
        .await?
        .ok_or_else(|| AppError::NotFound("review".to_string()))?;

    // レビューを更新
    review_repo
        .update(review_id, input.rating, input.comment.as_deref())
        .await?;

    Ok(())
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use std::sync::Mutex;

    use super::*;
    use crate::app::features::review::repository::{
        MovieExistsRepository, MovieRuntimeRow, PlatformExistsRepository, PointRepository,
        ReviewRepository, ReviewRow, WatchHistoryRepository,
    };
    use crate::app::responses::AppError;

    // ── モック ────────────────────────────────────────────────────

    struct MockReviewRepository {
        existing: Option<ReviewRow>, // find_by_movie_id / find_by_id が返す値
        created_id: i32,
        should_fail: bool,
    }

    impl MockReviewRepository {
        fn empty(next_id: i32) -> Self {
            Self { existing: None, created_id: next_id, should_fail: false }
        }
        fn with_existing(id: i32, user_id: i32, movie_id: i32) -> Self {
            Self {
                existing: Some(ReviewRow { id, user_id, movie_id }),
                created_id: id,
                should_fail: false,
            }
        }
        fn failing() -> Self {
            Self { existing: None, created_id: 0, should_fail: true }
        }
    }

    impl ReviewRepository for MockReviewRepository {
        async fn find_by_movie_id(&self, _user_id: i32, _movie_id: i32) -> Result<Option<ReviewRow>, AppError> {
            if self.should_fail { return Err(AppError::InternalServerError); }
            Ok(self.existing.as_ref().map(|r| ReviewRow { id: r.id, user_id: r.user_id, movie_id: r.movie_id }))
        }
        async fn find_by_id(&self, _user_id: i32, _review_id: i32) -> Result<Option<ReviewRow>, AppError> {
            if self.should_fail { return Err(AppError::InternalServerError); }
            Ok(self.existing.as_ref().map(|r| ReviewRow { id: r.id, user_id: r.user_id, movie_id: r.movie_id }))
        }
        async fn create(&self, _u: i32, _m: i32, _r: Option<f64>, _c: Option<&str>) -> Result<i32, AppError> {
            if self.should_fail { return Err(AppError::InternalServerError); }
            Ok(self.created_id)
        }
        async fn update(&self, _id: i32, _r: Option<f64>, _c: Option<&str>) -> Result<(), AppError> {
            if self.should_fail { return Err(AppError::InternalServerError); }
            Ok(())
        }
    }

    struct MockMovieRepo { runtime: Option<i32> }
    impl MockMovieRepo {
        fn with_runtime(rt: i32) -> Self { Self { runtime: Some(rt) } }
        fn not_found() -> Self { Self { runtime: None } }
    }
    impl MovieExistsRepository for MockMovieRepo {
        async fn find_by_id(&self, _id: i32) -> Result<Option<MovieRuntimeRow>, AppError> {
            Ok(self.runtime.map(|rt| MovieRuntimeRow { id: 1, runtime_minutes: rt }))
        }
    }

    struct MockPlatformRepo { exists: bool }
    impl MockPlatformRepo {
        fn found() -> Self { Self { exists: true } }
        fn not_found() -> Self { Self { exists: false } }
    }
    impl PlatformExistsRepository for MockPlatformRepo {
        async fn exists(&self, _id: i32) -> Result<bool, AppError> { Ok(self.exists) }
    }

    struct MockWatchHistoryRepo { next_id: i32 }
    impl MockWatchHistoryRepo { fn new(id: i32) -> Self { Self { next_id: id } } }
    impl WatchHistoryRepository for MockWatchHistoryRepo {
        async fn create(&self, _u: i32, _m: i32, _p: i32, _d: Option<&str>) -> Result<i32, AppError> {
            Ok(self.next_id)
        }
    }

    struct MockPointRepo { calls: Mutex<Vec<(i32, String)>> }
    impl MockPointRepo {
        fn new() -> Self { Self { calls: Mutex::new(vec![]) } }
        fn actions(&self) -> Vec<String> {
            self.calls.lock().unwrap().iter().map(|(_, a)| a.clone()).collect()
        }
    }
    impl PointRepository for MockPointRepo {
        async fn grant_points(&self, user_id: i32, points: i32, action: &str, _ref_id: i32) -> Result<(), AppError> {
            self.calls.lock().unwrap().push((points, action.to_string()));
            let _ = user_id;
            Ok(())
        }
    }

    // ── create_review テスト ──────────────────────────────────────

    #[tokio::test]
    async fn test_create_review_success_without_watch_history() {
        let review_repo = MockReviewRepository::empty(1);
        let movie_repo = MockMovieRepo::with_runtime(120);
        let platform_repo = MockPlatformRepo::found();
        let wh_repo = MockWatchHistoryRepo::new(1);
        let point_repo = MockPointRepo::new();

        let input = CreateReviewInput {
            rating: Some(4.0),
            comment: None,
            watch_history: None,
        };

        let result = create_review(&review_repo, &movie_repo, &platform_repo, &wh_repo, &point_repo, 1, 1, input).await;

        assert!(result.is_ok());
        // レビューポイントのみ付与
        assert_eq!(point_repo.actions(), vec!["review"]);
    }

    #[tokio::test]
    async fn test_create_review_success_with_watch_history() {
        let review_repo = MockReviewRepository::empty(1);
        let movie_repo = MockMovieRepo::with_runtime(120); // 120分 → 15pt
        let platform_repo = MockPlatformRepo::found();
        let wh_repo = MockWatchHistoryRepo::new(10);
        let point_repo = MockPointRepo::new();

        let input = CreateReviewInput {
            rating: Some(4.0),
            comment: None,
            watch_history: Some(WatchHistoryInput {
                platform_id: 1,
                watched_date: Some("2024-01-01".to_string()),
            }),
        };

        let result = create_review(&review_repo, &movie_repo, &platform_repo, &wh_repo, &point_repo, 1, 1, input).await;

        assert!(result.is_ok());
        // レビューポイントと視聴履歴ポイントの両方が付与される
        let actions = point_repo.actions();
        assert_eq!(actions.len(), 2);
        assert_eq!(actions[0], "review");
        assert_eq!(actions[1], "watch_history");
    }

    #[tokio::test]
    async fn test_create_review_fails_when_rating_and_comment_both_none() {
        let review_repo = MockReviewRepository::empty(1);
        let movie_repo = MockMovieRepo::with_runtime(120);
        let platform_repo = MockPlatformRepo::found();
        let wh_repo = MockWatchHistoryRepo::new(1);
        let point_repo = MockPointRepo::new();

        let input = CreateReviewInput { rating: None, comment: None, watch_history: None };

        let result = create_review(&review_repo, &movie_repo, &platform_repo, &wh_repo, &point_repo, 1, 1, input).await;

        assert!(matches!(result, Err(AppError::ValidationError(_))));
    }

    #[tokio::test]
    async fn test_create_review_fails_when_rating_out_of_range() {
        let review_repo = MockReviewRepository::empty(1);
        let movie_repo = MockMovieRepo::with_runtime(120);
        let platform_repo = MockPlatformRepo::found();
        let wh_repo = MockWatchHistoryRepo::new(1);
        let point_repo = MockPointRepo::new();

        let input = CreateReviewInput { rating: Some(5.1), comment: None, watch_history: None };

        let result = create_review(&review_repo, &movie_repo, &platform_repo, &wh_repo, &point_repo, 1, 1, input).await;

        assert!(matches!(result, Err(AppError::ValidationError(_))));
    }

    #[tokio::test]
    async fn test_create_review_fails_when_movie_not_found() {
        let review_repo = MockReviewRepository::empty(1);
        let movie_repo = MockMovieRepo::not_found();
        let platform_repo = MockPlatformRepo::found();
        let wh_repo = MockWatchHistoryRepo::new(1);
        let point_repo = MockPointRepo::new();

        let input = CreateReviewInput { rating: Some(4.0), comment: None, watch_history: None };

        let result = create_review(&review_repo, &movie_repo, &platform_repo, &wh_repo, &point_repo, 1, 1, input).await;

        assert!(matches!(result, Err(AppError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_create_review_fails_when_review_already_exists() {
        let review_repo = MockReviewRepository::with_existing(1, 1, 1);
        let movie_repo = MockMovieRepo::with_runtime(120);
        let platform_repo = MockPlatformRepo::found();
        let wh_repo = MockWatchHistoryRepo::new(1);
        let point_repo = MockPointRepo::new();

        let input = CreateReviewInput { rating: Some(4.0), comment: None, watch_history: None };

        let result = create_review(&review_repo, &movie_repo, &platform_repo, &wh_repo, &point_repo, 1, 1, input).await;

        assert!(matches!(result, Err(AppError::Conflict(_))));
    }

    #[tokio::test]
    async fn test_create_review_fails_when_platform_not_found() {
        let review_repo = MockReviewRepository::empty(1);
        let movie_repo = MockMovieRepo::with_runtime(120);
        let platform_repo = MockPlatformRepo::not_found();
        let wh_repo = MockWatchHistoryRepo::new(1);
        let point_repo = MockPointRepo::new();

        let input = CreateReviewInput {
            rating: Some(4.0),
            comment: None,
            watch_history: Some(WatchHistoryInput { platform_id: 99, watched_date: None }),
        };

        let result = create_review(&review_repo, &movie_repo, &platform_repo, &wh_repo, &point_repo, 1, 1, input).await;

        assert!(matches!(result, Err(AppError::NotFound(_))));
    }

    // ── calc_watch_points テスト ──────────────────────────────────

    #[test]
    fn test_calc_watch_points_short() {
        assert_eq!(calc_watch_points(90), 10);
        assert_eq!(calc_watch_points(60), 10);
    }

    #[test]
    fn test_calc_watch_points_medium() {
        assert_eq!(calc_watch_points(91), 15);
        assert_eq!(calc_watch_points(150), 15);
    }

    #[test]
    fn test_calc_watch_points_long() {
        assert_eq!(calc_watch_points(151), 20);
        assert_eq!(calc_watch_points(180), 20);
    }

    // ── update_review テスト ──────────────────────────────────────

    #[tokio::test]
    async fn test_update_review_success() {
        let review_repo = MockReviewRepository::with_existing(1, 1, 1);

        let input = UpdateReviewInput { rating: Some(5.0), comment: None };
        let result = update_review(&review_repo, 1, 1, input).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_update_review_success_comment_only() {
        let review_repo = MockReviewRepository::with_existing(1, 1, 1);

        let input = UpdateReviewInput { rating: None, comment: Some("updated".to_string()) };
        let result = update_review(&review_repo, 1, 1, input).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_update_review_fails_when_both_none() {
        let review_repo = MockReviewRepository::with_existing(1, 1, 1);

        let input = UpdateReviewInput { rating: None, comment: None };
        let result = update_review(&review_repo, 1, 1, input).await;

        assert!(matches!(result, Err(AppError::ValidationError(_))));
    }

    #[tokio::test]
    async fn test_update_review_fails_when_not_found() {
        let review_repo = MockReviewRepository::empty(0);

        let input = UpdateReviewInput { rating: Some(4.0), comment: None };
        let result = update_review(&review_repo, 1, 99, input).await;

        assert!(matches!(result, Err(AppError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_update_review_fails_when_rating_out_of_range() {
        let review_repo = MockReviewRepository::with_existing(1, 1, 1);

        let input = UpdateReviewInput { rating: Some(0.0), comment: None };
        let result = update_review(&review_repo, 1, 1, input).await;

        assert!(matches!(result, Err(AppError::ValidationError(_))));
    }

    #[tokio::test]
    async fn test_update_review_propagates_repo_error() {
        let review_repo = MockReviewRepository::failing();

        // find_by_id が失敗 → InternalServerError
        let input = UpdateReviewInput { rating: Some(4.0), comment: None };
        let result = update_review(&review_repo, 1, 1, input).await;

        assert!(matches!(result, Err(AppError::InternalServerError)));
    }
}
