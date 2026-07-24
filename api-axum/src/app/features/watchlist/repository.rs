use chrono::NaiveDate;
use sqlx::MySqlPool;
use tracing::error;

use crate::app::responses::AppError;

// ─── DB 行型 ──────────────────────────────────────────────────

/// watchlist + movies + GROUP_CONCAT でジャンルを結合した行
#[derive(Debug, sqlx::FromRow)]
pub struct WatchlistRow {
    pub id: i32,
    pub added_at: Option<chrono::DateTime<chrono::Utc>>,
    pub priority: Option<i32>,
    pub movie_id: i32,
    pub movie_title: String,
    pub movie_overview: String,
    pub movie_release_date: Option<NaiveDate>,
    pub movie_runtime_minutes: i32,
    pub movie_poster_url: Option<String>,
    pub movie_tmdb_id: i32,
    pub genre_codes: Option<String>,
    pub genre_names: Option<String>,
}

// ─── Repository trait ─────────────────────────────────────────

pub trait WatchlistRepository {
    /// ユーザーのウォッチリストを映画・ジャンル情報付きで取得する
    async fn find_by_user_id(
        &self,
        user_id: i32,
        limit: i32,
        offset: i32,
    ) -> Result<Vec<WatchlistRow>, AppError>;

    /// ウォッチリストに登録する
    async fn create(&self, user_id: i32, movie_id: i32, priority: i32) -> Result<(), AppError>;

    /// ウォッチリストIDに一致するレコードを削除する。削除件数を返す。
    async fn delete_by_id(&self, watchlist_id: i32) -> Result<u64, AppError>;
}

pub trait MovieExistsRepository {
    /// 映画の存在確認
    async fn exists(&self, movie_id: i32) -> Result<bool, AppError>;
}

// ─── MySQL 実装 ───────────────────────────────────────────────

pub struct MySqlWatchlistRepository<'a>(pub &'a MySqlPool);
pub struct MySqlMovieExistsRepository<'a>(pub &'a MySqlPool);

impl WatchlistRepository for MySqlWatchlistRepository<'_> {
    async fn find_by_user_id(
        &self,
        user_id: i32,
        limit: i32,
        offset: i32,
    ) -> Result<Vec<WatchlistRow>, AppError> {
        sqlx::query_as(
            "SELECT
                 w.id,
                 w.added_at,
                 w.priority,
                 m.id             AS movie_id,
                 m.title          AS movie_title,
                 m.overview       AS movie_overview,
                 m.release_date   AS movie_release_date,
                 m.runtime_minutes AS movie_runtime_minutes,
                 m.poster_url     AS movie_poster_url,
                 m.tmdb_id        AS movie_tmdb_id,
                 GROUP_CONCAT(DISTINCT g.code ORDER BY g.code) AS genre_codes,
                 GROUP_CONCAT(DISTINCT g.name ORDER BY g.code) AS genre_names
             FROM watchlist w
             INNER JOIN movies m ON m.id = w.movie_id
             LEFT  JOIN movie_genres mg ON mg.movie_id = m.id
             LEFT  JOIN genres g ON g.id = mg.genre_id
             WHERE w.user_id = ?
             GROUP BY w.id, m.id
             ORDER BY w.added_at DESC
             LIMIT ? OFFSET ?",
        )
        .bind(user_id)
        .bind(limit)
        .bind(offset)
        .fetch_all(self.0)
        .await
        .map_err(|e| {
            error!("find_by_user_id (watchlist) failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })
    }

    async fn create(&self, user_id: i32, movie_id: i32, priority: i32) -> Result<(), AppError> {
        sqlx::query(
            "INSERT INTO watchlist (user_id, movie_id, priority, added_at, created_at, updated_at)
             VALUES (?, ?, ?, NOW(), NOW(), NOW())",
        )
        .bind(user_id)
        .bind(movie_id)
        .bind(priority)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("watchlist create failed: {e}");
            AppError::InternalServerError
        })?;

        Ok(())
    }

    async fn delete_by_id(&self, watchlist_id: i32) -> Result<u64, AppError> {
        let result = sqlx::query("DELETE FROM watchlist WHERE id = ?")
            .bind(watchlist_id)
            .execute(self.0)
            .await
            .map_err(|e| {
                error!("watchlist delete_by_id failed (id={watchlist_id}): {e}");
                AppError::InternalServerError
            })?;

        Ok(result.rows_affected())
    }
}

impl MovieExistsRepository for MySqlMovieExistsRepository<'_> {
    async fn exists(&self, movie_id: i32) -> Result<bool, AppError> {
        let count: i64 = sqlx::query_scalar("SELECT COUNT(*) FROM movies WHERE id = ?")
            .bind(movie_id)
            .fetch_one(self.0)
            .await
            .map_err(|e| {
                error!("movie exists failed (movie_id={movie_id}): {e}");
                AppError::InternalServerError
            })?;

        Ok(count > 0)
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

    async fn insert_test_movie(pool: &MySqlPool, tmdb_id: i32) -> i32 {
        let result = sqlx::query(
            "INSERT INTO movies (tmdb_id, title, overview, release_date, runtime_minutes, poster_url, created_at, updated_at)
             VALUES (?, ?, '', '2024-01-01', 120, NULL, NOW(), NOW())",
        )
        .bind(tmdb_id)
        .bind(unique("movie"))
        .execute(pool)
        .await
        .expect("insert movie");
        result.last_insert_id() as i32
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_watchlist_create_and_find() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool).await;
        let movie_id = insert_test_movie(&pool, 99001).await;
        let repo = MySqlWatchlistRepository(&pool);

        repo.create(user_id, movie_id, 1).await.unwrap();

        let rows = repo.find_by_user_id(user_id, 12, 0).await.unwrap();
        assert_eq!(rows.len(), 1);
        assert_eq!(rows[0].movie_id, movie_id);
        assert_eq!(rows[0].priority, Some(1));
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_watchlist_find_returns_empty_for_unknown_user() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlWatchlistRepository(&pool);

        let rows = repo.find_by_user_id(999999, 12, 0).await.unwrap();
        assert!(rows.is_empty());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_watchlist_delete_returns_affected_count() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool).await;
        let movie_id = insert_test_movie(&pool, 99002).await;
        let repo = MySqlWatchlistRepository(&pool);

        repo.create(user_id, movie_id, 1).await.unwrap();

        let id: i32 = sqlx::query_scalar("SELECT id FROM watchlist WHERE user_id = ? ORDER BY id DESC LIMIT 1")
            .bind(user_id)
            .fetch_one(&pool)
            .await
            .unwrap();

        let affected = repo.delete_by_id(id).await.unwrap();
        assert_eq!(affected, 1);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_watchlist_delete_returns_zero_for_unknown_id() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlWatchlistRepository(&pool);

        let affected = repo.delete_by_id(999999).await.unwrap();
        assert_eq!(affected, 0);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_movie_exists_returns_false_for_unknown() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMovieExistsRepository(&pool);

        let result = repo.exists(999999).await.unwrap();
        assert!(!result);
    }
}
