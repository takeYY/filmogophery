use sqlx::MySqlPool;
use tracing::error;

use crate::app::responses::AppError;

// ─── DB 行型 ──────────────────────────────────────────────────

#[derive(Debug, sqlx::FromRow)]
pub struct ReviewRow {
    #[allow(dead_code)]
    pub id: i32,
    #[allow(dead_code)]
    pub user_id: i32,
    #[allow(dead_code)]
    pub movie_id: i32,
}

#[derive(Debug, sqlx::FromRow)]
pub struct MovieRuntimeRow {
    #[allow(dead_code)]
    pub id: i32,
    pub runtime_minutes: i32,
}

// ─── Repository trait ─────────────────────────────────────────

pub trait ReviewRepository {
    /// ユーザーの映画レビューが既に存在するか確認する
    async fn find_by_movie_id(
        &self,
        user_id: i32,
        movie_id: i32,
    ) -> Result<Option<ReviewRow>, AppError>;

    /// レビューIDとユーザーIDに一致するレビューを取得する（所有権チェック込み）
    async fn find_by_id(
        &self,
        user_id: i32,
        review_id: i32,
    ) -> Result<Option<ReviewRow>, AppError>;

    /// レビューを作成し、採番された ID を返す
    async fn create(
        &self,
        user_id: i32,
        movie_id: i32,
        rating: Option<f64>,
        comment: Option<&str>,
    ) -> Result<i32, AppError>;

    /// レビューを更新する
    async fn update(
        &self,
        review_id: i32,
        rating: Option<f64>,
        comment: Option<&str>,
    ) -> Result<(), AppError>;
}

pub trait MovieExistsRepository {
    /// 映画の存在確認（runtime_minutes も取得）
    async fn find_by_id(&self, movie_id: i32) -> Result<Option<MovieRuntimeRow>, AppError>;
}

pub trait PlatformExistsRepository {
    /// プラットフォームの存在確認
    async fn exists(&self, platform_id: i32) -> Result<bool, AppError>;
}

pub trait WatchHistoryRepository {
    /// 視聴履歴を作成し、採番された ID を返す
    async fn create(
        &self,
        user_id: i32,
        movie_id: i32,
        platform_id: i32,
        watched_date: Option<&str>,
    ) -> Result<i32, AppError>;
}

pub trait PointRepository {
    /// user_points を upsert し、ポイント履歴を記録する（トランザクション外）
    async fn grant_points(
        &self,
        user_id: i32,
        points: i32,
        action: &str,
        reference_id: i32,
    ) -> Result<(), AppError>;
}

// ─── MySQL 実装 ───────────────────────────────────────────────

pub struct MySqlReviewRepository<'a>(pub &'a MySqlPool);
pub struct MySqlMovieExistsRepository<'a>(pub &'a MySqlPool);
pub struct MySqlPlatformExistsRepository<'a>(pub &'a MySqlPool);
pub struct MySqlWatchHistoryRepository<'a>(pub &'a MySqlPool);
pub struct MySqlPointRepository<'a>(pub &'a MySqlPool);

impl ReviewRepository for MySqlReviewRepository<'_> {
    async fn find_by_movie_id(
        &self,
        user_id: i32,
        movie_id: i32,
    ) -> Result<Option<ReviewRow>, AppError> {
        sqlx::query_as(
            "SELECT id, user_id, movie_id
             FROM reviews
             WHERE user_id = ? AND movie_id = ?
             LIMIT 1",
        )
        .bind(user_id)
        .bind(movie_id)
        .fetch_optional(self.0)
        .await
        .map_err(|e| {
            error!("find_by_movie_id failed: {e}");
            AppError::InternalServerError
        })
    }

    async fn find_by_id(
        &self,
        user_id: i32,
        review_id: i32,
    ) -> Result<Option<ReviewRow>, AppError> {
        sqlx::query_as(
            "SELECT id, user_id, movie_id
             FROM reviews
             WHERE id = ? AND user_id = ?
             LIMIT 1",
        )
        .bind(review_id)
        .bind(user_id)
        .fetch_optional(self.0)
        .await
        .map_err(|e| {
            error!("find_by_id failed (review_id={review_id}): {e}");
            AppError::InternalServerError
        })
    }

    async fn create(
        &self,
        user_id: i32,
        movie_id: i32,
        rating: Option<f64>,
        comment: Option<&str>,
    ) -> Result<i32, AppError> {
        let result = sqlx::query(
            "INSERT INTO reviews (user_id, movie_id, rating, comment, created_at, updated_at)
             VALUES (?, ?, ?, ?, NOW(), NOW())",
        )
        .bind(user_id)
        .bind(movie_id)
        .bind(rating)
        .bind(comment)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("review create failed: {e}");
            if e.to_string().contains("1062") || e.to_string().contains("Duplicate entry") {
                return AppError::Conflict("review".to_string());
            }
            AppError::InternalServerError
        })?;

        Ok(result.last_insert_id() as i32)
    }

    async fn update(
        &self,
        review_id: i32,
        rating: Option<f64>,
        comment: Option<&str>,
    ) -> Result<(), AppError> {
        sqlx::query(
            "UPDATE reviews
             SET rating = COALESCE(?, rating),
                 comment = COALESCE(?, comment),
                 updated_at = NOW()
             WHERE id = ?",
        )
        .bind(rating)
        .bind(comment)
        .bind(review_id)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("review update failed (review_id={review_id}): {e}");
            AppError::InternalServerError
        })?;

        Ok(())
    }
}

impl MovieExistsRepository for MySqlMovieExistsRepository<'_> {
    async fn find_by_id(&self, movie_id: i32) -> Result<Option<MovieRuntimeRow>, AppError> {
        sqlx::query_as(
            "SELECT id, runtime_minutes FROM movies WHERE id = ? LIMIT 1",
        )
        .bind(movie_id)
        .fetch_optional(self.0)
        .await
        .map_err(|e| {
            error!("movie find_by_id failed (movie_id={movie_id}): {e}");
            AppError::InternalServerError
        })
    }
}

impl PlatformExistsRepository for MySqlPlatformExistsRepository<'_> {
    async fn exists(&self, platform_id: i32) -> Result<bool, AppError> {
        let count: i64 = sqlx::query_scalar(
            "SELECT COUNT(*) FROM platforms WHERE id = ?",
        )
        .bind(platform_id)
        .fetch_one(self.0)
        .await
        .map_err(|e| {
            error!("platform exists failed (platform_id={platform_id}): {e}");
            AppError::InternalServerError
        })?;

        Ok(count > 0)
    }
}

impl WatchHistoryRepository for MySqlWatchHistoryRepository<'_> {
    async fn create(
        &self,
        user_id: i32,
        movie_id: i32,
        platform_id: i32,
        watched_date: Option<&str>,
    ) -> Result<i32, AppError> {
        let result = sqlx::query(
            "INSERT INTO watch_history (user_id, movie_id, platform_id, watched_date, created_at, updated_at)
             VALUES (?, ?, ?, ?, NOW(), NOW())",
        )
        .bind(user_id)
        .bind(movie_id)
        .bind(platform_id)
        .bind(watched_date)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("watch_history create failed: {e}");
            AppError::InternalServerError
        })?;

        Ok(result.last_insert_id() as i32)
    }
}

impl PointRepository for MySqlPointRepository<'_> {
    async fn grant_points(
        &self,
        user_id: i32,
        points: i32,
        action: &str,
        reference_id: i32,
    ) -> Result<(), AppError> {
        // user_points を upsert
        sqlx::query(
            "INSERT INTO user_points (user_id, total_points, level, created_at, updated_at)
             VALUES (?, ?, 1, NOW(), NOW())
             ON DUPLICATE KEY UPDATE
               total_points = total_points + VALUES(total_points),
               updated_at = NOW()",
        )
        .bind(user_id)
        .bind(points)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("grant_points upsert failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })?;

        // point_history に記録
        sqlx::query(
            "INSERT INTO point_history (user_id, points, action, reference_id, created_at)
             VALUES (?, ?, ?, ?, NOW())",
        )
        .bind(user_id)
        .bind(points)
        .bind(action)
        .bind(reference_id)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("grant_points history failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })?;

        Ok(())
    }
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use sqlx::MySqlPool;
    use uuid::Uuid;

    use super::*;

    async fn test_pool() -> Option<MySqlPool> {
        let url = std::env::var("DATABASE_URL").ok()?;
        MySqlPool::connect(&url).await.ok()
    }

    fn unique(base: &str) -> String {
        format!("{}_{}", base, Uuid::new_v4().simple())
    }

    async fn insert_test_user(pool: &MySqlPool) -> i32 {
        let result = sqlx::query(
            "INSERT INTO users (username, email, password_hash, is_active, last_login_at, created_at, updated_at)
             VALUES (?, ?, 'hash', 1, NOW(), NOW(), NOW())",
        )
        .bind(unique("user"))
        .bind(unique("user@example.com"))
        .execute(pool)
        .await
        .expect("insert user");
        result.last_insert_id() as i32
    }

    async fn insert_test_movie(pool: &MySqlPool) -> i32 {
        let result = sqlx::query(
            "INSERT INTO movies (tmdb_id, title, overview, release_date, runtime_minutes, poster_url, created_at, updated_at)
             VALUES (?, ?, '', '2024-01-01', 120, NULL, NOW(), NOW())",
        )
        .bind(rand_i32())
        .bind(unique("movie"))
        .execute(pool)
        .await
        .expect("insert movie");
        result.last_insert_id() as i32
    }

    fn rand_i32() -> i32 {
        use std::collections::hash_map::DefaultHasher;
        use std::hash::{Hash, Hasher};
        let mut h = DefaultHasher::new();
        Uuid::new_v4().hash(&mut h);
        (h.finish() % 1_000_000) as i32 + 1
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_review_create_and_find() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool).await;
        let movie_id = insert_test_movie(&pool).await;
        let repo = MySqlReviewRepository(&pool);

        // 作成
        let review_id = repo.create(user_id, movie_id, Some(4.0), Some("good")).await;
        assert!(review_id.is_ok());
        let review_id = review_id.unwrap();
        assert!(review_id > 0);

        // 映画ID で検索
        let found = repo.find_by_movie_id(user_id, movie_id).await.unwrap();
        assert!(found.is_some());
        assert_eq!(found.unwrap().id, review_id);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_review_create_duplicate_returns_conflict() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool).await;
        let movie_id = insert_test_movie(&pool).await;
        let repo = MySqlReviewRepository(&pool);

        repo.create(user_id, movie_id, Some(4.0), None).await.unwrap();
        let result = repo.create(user_id, movie_id, Some(3.0), None).await;

        assert!(matches!(result, Err(AppError::Conflict(_))));
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_review_find_by_id_returns_none_for_wrong_user() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool).await;
        let movie_id = insert_test_movie(&pool).await;
        let repo = MySqlReviewRepository(&pool);

        let review_id = repo.create(user_id, movie_id, Some(4.0), None).await.unwrap();

        // 別ユーザーからは取得できない
        let result = repo.find_by_id(user_id + 1, review_id).await.unwrap();
        assert!(result.is_none());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_review_update_succeeds() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool).await;
        let movie_id = insert_test_movie(&pool).await;
        let repo = MySqlReviewRepository(&pool);

        let review_id = repo.create(user_id, movie_id, Some(3.0), None).await.unwrap();
        let result = repo.update(review_id, Some(5.0), Some("updated")).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_movie_find_by_id_returns_none_for_unknown() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMovieExistsRepository(&pool);

        let result = repo.find_by_id(999999).await.unwrap();
        assert!(result.is_none());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_platform_exists_returns_false_for_unknown() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlPlatformExistsRepository(&pool);

        let result = repo.exists(999999).await.unwrap();
        assert!(!result);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_grant_points_upserts_and_records_history() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool).await;
        let repo = MySqlPointRepository(&pool);

        let result = repo.grant_points(user_id, 20, "review", 1).await;
        assert!(result.is_ok());

        let total: i32 = sqlx::query_scalar(
            "SELECT total_points FROM user_points WHERE user_id = ?",
        )
        .bind(user_id)
        .fetch_one(&pool)
        .await
        .unwrap();
        assert_eq!(total, 20);

        // 2回目: 加算される
        repo.grant_points(user_id, 10, "watch_history", 2).await.unwrap();
        let total2: i32 = sqlx::query_scalar(
            "SELECT total_points FROM user_points WHERE user_id = ?",
        )
        .bind(user_id)
        .fetch_one(&pool)
        .await
        .unwrap();
        assert_eq!(total2, 30);
    }
}
