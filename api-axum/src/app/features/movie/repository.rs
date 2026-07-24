use chrono::NaiveDate;
use sqlx::MySqlPool;
use tracing::error;

use crate::app::responses::AppError;

// ─── DB 行型 ──────────────────────────────────────────────────

/// movies テーブル + GROUP_CONCAT でジャンルを結合した行
#[derive(Debug, sqlx::FromRow)]
pub struct MovieRow {
    pub id: i32,
    pub title: String,
    pub overview: String,
    pub release_date: Option<NaiveDate>,
    pub runtime_minutes: i32,
    pub poster_url: Option<String>,
    pub tmdb_id: i32,
    /// "ACTION,DRAMA" 形式（NULL の場合は空文字列）
    pub genre_codes: Option<String>,
    /// "アクション,ドラマ" 形式
    pub genre_names: Option<String>,
}

/// movies テーブル + reviews join (詳細用)
#[derive(Debug, sqlx::FromRow)]
pub struct MovieDetailRow {
    pub id: i32,
    pub title: String,
    pub overview: String,
    pub release_date: Option<NaiveDate>,
    pub runtime_minutes: i32,
    pub poster_url: Option<String>,
    pub tmdb_id: i32,
    pub genre_codes: Option<String>,
    pub genre_names: Option<String>,
    // review (nullable)
    pub review_id: Option<i32>,
    pub review_rating: Option<f64>,
    pub review_comment: Option<String>,
    pub review_created_at: Option<chrono::DateTime<chrono::Utc>>,
    pub review_updated_at: Option<chrono::DateTime<chrono::Utc>>,
    // series (nullable)
    pub series_name: Option<String>,
    pub series_poster_url: Option<String>,
}

/// tmdb_id → DB id のマッピング（search 用）
#[derive(Debug, sqlx::FromRow)]
pub struct MovieTmdbRow {
    pub id: i32,
    pub tmdb_id: i32,
    pub genre_codes: Option<String>,
    pub genre_names: Option<String>,
}

/// watch_history + platform + movie + genres 結合行（ユーザー視聴履歴一覧用）
#[derive(Debug, sqlx::FromRow)]
pub struct WatchHistoryRow {
    pub id: i32,
    pub watched_date: Option<chrono::NaiveDate>,
    // platform
    pub platform_id: i32,
    pub platform_code: String,
    pub platform_name: String,
    // movie
    pub movie_id: i32,
    pub movie_title: String,
    pub movie_overview: String,
    pub movie_release_date: Option<chrono::NaiveDate>,
    pub movie_runtime_minutes: i32,
    pub movie_poster_url: Option<String>,
    pub movie_tmdb_id: i32,
    pub genre_codes: Option<String>,
    pub genre_names: Option<String>,
}

/// watch_history + platform 結合行（映画ごとの視聴履歴用）
#[derive(Debug, sqlx::FromRow)]
pub struct MovieWatchHistoryRow {
    pub id: i32,
    pub watched_date: Option<chrono::NaiveDate>,
    pub platform_id: i32,
    pub platform_code: String,
    pub platform_name: String,
}

// ─── Repository trait ─────────────────────────────────────────

pub trait MovieRepository {
    /// ログインユーザーがレビューした映画一覧（ジャンル絞り込み可）
    async fn find_reviewed_by_user(
        &self,
        user_id: i32,
        genre: Option<&str>,
        limit: i32,
        offset: i32,
    ) -> Result<Vec<MovieRow>, AppError>;

    /// ID で映画詳細を取得（review・series を LEFT JOIN）
    async fn find_detail_by_id(
        &self,
        movie_id: i32,
        user_id: i32,
    ) -> Result<Option<MovieDetailRow>, AppError>;

    /// tmdb_id リストで映画を取得
    async fn find_by_tmdb_ids(&self, tmdb_ids: &[i32]) -> Result<Vec<MovieTmdbRow>, AppError>;

    /// 上映時間を更新
    async fn update_runtime_minutes(&self, movie_id: i32, runtime: i32) -> Result<(), AppError>;

    /// 映画を一括挿入（search で新規映画をDBに登録）
    async fn batch_insert(
        &self,
        movies: &[NewMovieInput],
    ) -> Result<(), AppError>;

    /// 映画の視聴履歴一覧を取得（プラットフォーム情報付き）
    async fn find_watch_history_by_movie_id(
        &self,
        user_id: i32,
        movie_id: i32,
    ) -> Result<Vec<MovieWatchHistoryRow>, AppError>;

    /// ユーザーの視聴履歴一覧を取得（映画・プラットフォーム・ジャンル情報付き）
    async fn find_watch_history_by_user_id(
        &self,
        user_id: i32,
        limit: i32,
        offset: i32,
    ) -> Result<Vec<WatchHistoryRow>, AppError>;

    /// 映画の存在確認（軽量）
    async fn movie_exists(&self, movie_id: i32) -> Result<bool, AppError>;
}

pub struct NewMovieInput {
    pub tmdb_id: i32,
    pub title: String,
    pub overview: String,
    pub release_date: String,   // "YYYY-MM-DD"
    pub poster_url: Option<String>,
    pub genre_ids: Vec<i32>,
}

// ─── MySQL 実装 ───────────────────────────────────────────────

pub struct MySqlMovieRepository<'a>(pub &'a MySqlPool);

impl MovieRepository for MySqlMovieRepository<'_> {
    async fn find_reviewed_by_user(
        &self,
        user_id: i32,
        genre: Option<&str>,
        limit: i32,
        offset: i32,
    ) -> Result<Vec<MovieRow>, AppError> {
        // ジャンル指定あり・なしで HAVING 句を切り替える
        let rows: Vec<MovieRow> = if let Some(g) = genre.filter(|s| !s.is_empty()) {
            sqlx::query_as(
                "SELECT
                     m.id, m.title, m.overview, m.release_date,
                     m.runtime_minutes, m.poster_url, m.tmdb_id,
                     GROUP_CONCAT(DISTINCT g.code ORDER BY g.code) AS genre_codes,
                     GROUP_CONCAT(DISTINCT g.name ORDER BY g.code) AS genre_names
                 FROM movies m
                 INNER JOIN reviews r ON r.movie_id = m.id AND r.user_id = ?
                 LEFT  JOIN movie_genres mg ON mg.movie_id = m.id
                 LEFT  JOIN genres g ON g.id = mg.genre_id
                 GROUP BY m.id
                 HAVING FIND_IN_SET(?, GROUP_CONCAT(DISTINCT g.code))
                 ORDER BY MAX(r.created_at) DESC
                 LIMIT ? OFFSET ?",
            )
            .bind(user_id)
            .bind(g)
            .bind(limit)
            .bind(offset)
            .fetch_all(self.0)
            .await
        } else {
            sqlx::query_as(
                "SELECT
                     m.id, m.title, m.overview, m.release_date,
                     m.runtime_minutes, m.poster_url, m.tmdb_id,
                     GROUP_CONCAT(DISTINCT g.code ORDER BY g.code) AS genre_codes,
                     GROUP_CONCAT(DISTINCT g.name ORDER BY g.code) AS genre_names
                 FROM movies m
                 INNER JOIN reviews r ON r.movie_id = m.id AND r.user_id = ?
                 LEFT  JOIN movie_genres mg ON mg.movie_id = m.id
                 LEFT  JOIN genres g ON g.id = mg.genre_id
                 GROUP BY m.id
                 ORDER BY MAX(r.created_at) DESC
                 LIMIT ? OFFSET ?",
            )
            .bind(user_id)
            .bind(limit)
            .bind(offset)
            .fetch_all(self.0)
            .await
        }
        .map_err(|e| {
            error!("find_reviewed_by_user failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })?;

        Ok(rows)
    }

    async fn find_detail_by_id(
        &self,
        movie_id: i32,
        user_id: i32,
    ) -> Result<Option<MovieDetailRow>, AppError> {
        sqlx::query_as(
            "SELECT
                 m.id, m.title, m.overview, m.release_date,
                 m.runtime_minutes, m.poster_url, m.tmdb_id,
                 GROUP_CONCAT(DISTINCT g.code ORDER BY g.code) AS genre_codes,
                 GROUP_CONCAT(DISTINCT g.name ORDER BY g.code) AS genre_names,
                 r.id         AS review_id,
                 r.rating     AS review_rating,
                 r.comment    AS review_comment,
                 r.created_at AS review_created_at,
                 r.updated_at AS review_updated_at,
                 s.name       AS series_name,
                 s.poster_url AS series_poster_url
             FROM movies m
             LEFT JOIN reviews r ON r.movie_id = m.id AND r.user_id = ?
             LEFT JOIN movie_genres mg ON mg.movie_id = m.id
             LEFT JOIN genres g ON g.id = mg.genre_id
             LEFT JOIN series s ON s.id = m.series_id
             WHERE m.id = ?
             GROUP BY m.id, r.id, s.id",
        )
        .bind(user_id)
        .bind(movie_id)
        .fetch_optional(self.0)
        .await
        .map_err(|e| {
            error!("find_detail_by_id failed (movie_id={movie_id}): {e}");
            AppError::InternalServerError
        })
    }

    async fn find_by_tmdb_ids(&self, tmdb_ids: &[i32]) -> Result<Vec<MovieTmdbRow>, AppError> {
        if tmdb_ids.is_empty() {
            return Ok(vec![]);
        }

        // IN 句は動的に組み立てる
        let placeholders = tmdb_ids.iter().map(|_| "?").collect::<Vec<_>>().join(", ");
        let sql = format!(
            "SELECT
                 m.id, m.tmdb_id,
                 GROUP_CONCAT(DISTINCT g.code ORDER BY g.code) AS genre_codes,
                 GROUP_CONCAT(DISTINCT g.name ORDER BY g.code) AS genre_names
             FROM movies m
             LEFT JOIN movie_genres mg ON mg.movie_id = m.id
             LEFT JOIN genres g ON g.id = mg.genre_id
             WHERE m.tmdb_id IN ({placeholders})
             GROUP BY m.id"
        );

        let mut q = sqlx::query_as(&sql);
        for id in tmdb_ids {
            q = q.bind(id);
        }

        q.fetch_all(self.0).await.map_err(|e| {
            error!("find_by_tmdb_ids failed: {e}");
            AppError::InternalServerError
        })
    }

    async fn update_runtime_minutes(&self, movie_id: i32, runtime: i32) -> Result<(), AppError> {
        sqlx::query("UPDATE movies SET runtime_minutes = ? WHERE id = ?")
            .bind(runtime)
            .bind(movie_id)
            .execute(self.0)
            .await
            .map_err(|e| {
                error!("update_runtime_minutes failed (movie_id={movie_id}): {e}");
                AppError::InternalServerError
            })?;
        Ok(())
    }

    async fn batch_insert(&self, movies: &[NewMovieInput]) -> Result<(), AppError> {
        if movies.is_empty() {
            return Ok(());
        }

        for m in movies {
            // movies を INSERT（重複は IGNORE）
            let result = sqlx::query(
                "INSERT IGNORE INTO movies
                     (tmdb_id, title, overview, release_date, runtime_minutes, poster_url,
                      created_at, updated_at)
                 VALUES (?, ?, ?, ?, 0, ?, NOW(), NOW())",
            )
            .bind(m.tmdb_id)
            .bind(&m.title)
            .bind(&m.overview)
            .bind(&m.release_date)
            .bind(&m.poster_url)
            .execute(self.0)
            .await
            .map_err(|e| {
                error!("batch_insert movies failed: {e}");
                AppError::InternalServerError
            })?;

            if result.rows_affected() == 0 {
                // 重複で INSERT されなかった場合はジャンル紐付けをスキップ
                continue;
            }

            let movie_id = result.last_insert_id() as i32;

            // movie_genres を INSERT
            for &genre_id in &m.genre_ids {
                sqlx::query(
                    "INSERT IGNORE INTO movie_genres (movie_id, genre_id) VALUES (?, ?)",
                )
                .bind(movie_id)
                .bind(genre_id)
                .execute(self.0)
                .await
                .map_err(|e| {
                    error!("batch_insert movie_genres failed: {e}");
                    AppError::InternalServerError
                })?;
            }
        }

        Ok(())
    }

    async fn find_watch_history_by_movie_id(
        &self,
        user_id: i32,
        movie_id: i32,
    ) -> Result<Vec<MovieWatchHistoryRow>, AppError> {
        sqlx::query_as(
            "SELECT
                 wh.id,
                 wh.watched_date,
                 p.id   AS platform_id,
                 p.code AS platform_code,
                 p.name AS platform_name
             FROM watch_history wh
             INNER JOIN platforms p ON p.id = wh.platform_id
             WHERE wh.user_id = ? AND wh.movie_id = ?
             ORDER BY wh.watched_date DESC",
        )
        .bind(user_id)
        .bind(movie_id)
        .fetch_all(self.0)
        .await
        .map_err(|e| {
            error!("find_watch_history_by_movie_id failed: {e}");
            AppError::InternalServerError
        })
    }

    async fn find_watch_history_by_user_id(
        &self,
        user_id: i32,
        limit: i32,
        offset: i32,
    ) -> Result<Vec<WatchHistoryRow>, AppError> {
        sqlx::query_as(
            "SELECT
                 wh.id,
                 wh.watched_date,
                 p.id              AS platform_id,
                 p.code            AS platform_code,
                 p.name            AS platform_name,
                 m.id              AS movie_id,
                 m.title           AS movie_title,
                 m.overview        AS movie_overview,
                 m.release_date    AS movie_release_date,
                 m.runtime_minutes AS movie_runtime_minutes,
                 m.poster_url      AS movie_poster_url,
                 m.tmdb_id         AS movie_tmdb_id,
                 GROUP_CONCAT(DISTINCT g.code ORDER BY g.code) AS genre_codes,
                 GROUP_CONCAT(DISTINCT g.name ORDER BY g.code) AS genre_names
             FROM watch_history wh
             INNER JOIN platforms p ON p.id = wh.platform_id
             INNER JOIN movies m ON m.id = wh.movie_id
             LEFT  JOIN movie_genres mg ON mg.movie_id = m.id
             LEFT  JOIN genres g ON g.id = mg.genre_id
             WHERE wh.user_id = ?
             GROUP BY wh.id, p.id, m.id
             ORDER BY wh.watched_date DESC
             LIMIT ? OFFSET ?",
        )
        .bind(user_id)
        .bind(limit)
        .bind(offset)
        .fetch_all(self.0)
        .await
        .map_err(|e| {
            error!("find_watch_history_by_user_id failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })
    }

    async fn movie_exists(&self, movie_id: i32) -> Result<bool, AppError> {
        let count: i64 = sqlx::query_scalar("SELECT COUNT(*) FROM movies WHERE id = ?")
            .bind(movie_id)
            .fetch_one(self.0)
            .await
            .map_err(|e| {
                error!("movie_exists failed (movie_id={movie_id}): {e}");
                AppError::InternalServerError
            })?;
        Ok(count > 0)
    }
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use sqlx::MySqlPool;

    use super::*;

    async fn test_pool() -> Option<MySqlPool> {
        let url = std::env::var("DATABASE_URL").ok()?;
        MySqlPool::connect(&url).await.ok()
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_find_reviewed_by_user_returns_rows() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMovieRepository(&pool);

        // user_id=1 でレビュー済み映画を取得（データがなくても空リストで正常終了）
        let result = repo.find_reviewed_by_user(1, None, 12, 0).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_find_reviewed_by_user_with_genre_filter() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMovieRepository(&pool);

        let result = repo.find_reviewed_by_user(1, Some("ACTION"), 12, 0).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_find_detail_by_id_returns_none_for_unknown() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMovieRepository(&pool);

        let result = repo.find_detail_by_id(999999, 1).await;
        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_find_by_tmdb_ids_returns_empty_for_empty_input() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMovieRepository(&pool);

        let result = repo.find_by_tmdb_ids(&[]).await;
        assert!(result.is_ok());
        assert!(result.unwrap().is_empty());
    }
}
