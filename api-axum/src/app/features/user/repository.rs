use chrono::{DateTime, Utc};
use sqlx::MySqlPool;
use tracing::error;

use crate::app::responses::AppError;

// ─── DB 行型 ──────────────────────────────────────────────────

#[derive(Debug, Clone, sqlx::FromRow)]
pub struct UserRow {
    pub id: i32,
    pub username: String,
    /// ログイン実装時に使用
    #[allow(dead_code)]
    pub password_hash: String,
}

#[derive(Debug, Clone, sqlx::FromRow)]
pub struct UserPointsRow {
    pub total_points: i32,
    pub level: i32,
}

#[derive(Debug, Clone, sqlx::FromRow)]
pub struct PointHistoryRow {
    pub id: i32,
    pub points: i32,
    pub action: String,
    pub reference_id: i32,
    pub created_at: Option<DateTime<Utc>>,
}

// ─── Repository trait ────────────────────────────────────────

pub trait UserRepository {
    /// メールアドレスでユーザーを取得
    async fn find_by_email(&self, email: &str) -> Result<Option<UserRow>, AppError>;

    /// ID でアクティブなユーザーを取得
    async fn find_by_id(&self, id: i32) -> Result<Option<UserRow>, AppError>;

    /// ユーザーを作成し、採番された ID を返す
    async fn create(
        &self,
        username: &str,
        email: &str,
        password_hash: &str,
        now: DateTime<Utc>,
    ) -> Result<i32, AppError>;
}

pub trait PointRepository {
    /// ユーザーのポイントを取得（レコードがなければ None）
    async fn find_by_user_id(&self, user_id: i32) -> Result<Option<UserPointsRow>, AppError>;

    /// ポイント履歴を取得（降順）
    async fn find_history_by_user_id(
        &self,
        user_id: i32,
        limit: i32,
        offset: i32,
    ) -> Result<Vec<PointHistoryRow>, AppError>;
}

pub trait TokenRepository {
    /// リフレッシュトークンを保存する
    async fn save_refresh_token(
        &self,
        user_id: i32,
        token_hash: &str,
        expires_at: DateTime<Utc>,
        now: DateTime<Utc>,
    ) -> Result<(), AppError>;

    /// ユーザーの有効なリフレッシュトークンをすべて無効化する
    async fn revoke_active_tokens(
        &self,
        user_id: i32,
        now: DateTime<Utc>,
    ) -> Result<(), AppError>;
}

// ─── MySQL 実装 ───────────────────────────────────────────────

pub struct MySqlUserRepository<'a>(pub &'a MySqlPool);
pub struct MySqlPointRepository<'a>(pub &'a MySqlPool);
pub struct MySqlTokenRepository<'a>(pub &'a MySqlPool);

impl UserRepository for MySqlUserRepository<'_> {
    async fn find_by_email(&self, email: &str) -> Result<Option<UserRow>, AppError> {
        sqlx::query_as(
            "SELECT id, username, password_hash
             FROM users
             WHERE email = ? AND is_active = 1
             LIMIT 1",
        )
        .bind(email)
        .fetch_optional(self.0)
        .await
        .map_err(|e| {
            error!("find_by_email failed: {e}");
            AppError::InternalServerError
        })
    }

    async fn find_by_id(&self, id: i32) -> Result<Option<UserRow>, AppError> {
        sqlx::query_as(
            "SELECT id, username, password_hash
             FROM users
             WHERE id = ? AND is_active = 1
             LIMIT 1",
        )
        .bind(id)
        .fetch_optional(self.0)
        .await
        .map_err(|e| {
            error!("find_by_id failed (id={id}): {e}");
            AppError::InternalServerError
        })
    }

    async fn create(
        &self,
        username: &str,
        email: &str,
        password_hash: &str,
        now: DateTime<Utc>,
    ) -> Result<i32, AppError> {
        let result = sqlx::query(
            "INSERT INTO users
                 (username, email, password_hash, is_active, last_login_at, created_at, updated_at)
             VALUES (?, ?, ?, 1, ?, ?, ?)",
        )
        .bind(username)
        .bind(email)
        .bind(password_hash)
        .bind(now)
        .bind(now)
        .bind(now)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("user create failed: {e}");
            // MySQL の重複エラー（1062）を Conflict に変換
            if e.to_string().contains("1062") || e.to_string().contains("Duplicate entry") {
                return AppError::Conflict("user".to_string());
            }
            AppError::InternalServerError
        })?;

        Ok(result.last_insert_id() as i32)
    }
}

impl PointRepository for MySqlPointRepository<'_> {
    async fn find_by_user_id(&self, user_id: i32) -> Result<Option<UserPointsRow>, AppError> {
        sqlx::query_as(
            "SELECT total_points, level
             FROM user_points
             WHERE user_id = ?
             LIMIT 1",
        )
        .bind(user_id)
        .fetch_optional(self.0)
        .await
        .map_err(|e| {
            error!("find_by_user_id (points) failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })
    }

    async fn find_history_by_user_id(
        &self,
        user_id: i32,
        limit: i32,
        offset: i32,
    ) -> Result<Vec<PointHistoryRow>, AppError> {
        sqlx::query_as(
            "SELECT id, points, action, reference_id, created_at
             FROM point_history
             WHERE user_id = ?
             ORDER BY created_at DESC
             LIMIT ? OFFSET ?",
        )
        .bind(user_id)
        .bind(limit)
        .bind(offset)
        .fetch_all(self.0)
        .await
        .map_err(|e| {
            error!("find_history_by_user_id failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })
    }
}

impl TokenRepository for MySqlTokenRepository<'_> {
    async fn save_refresh_token(
        &self,
        user_id: i32,
        token_hash: &str,
        expires_at: DateTime<Utc>,
        now: DateTime<Utc>,
    ) -> Result<(), AppError> {
        sqlx::query(
            "INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at)
             VALUES (?, ?, ?, ?)",
        )
        .bind(user_id)
        .bind(token_hash)
        .bind(expires_at)
        .bind(now)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("save_refresh_token failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })?;

        Ok(())
    }

    async fn revoke_active_tokens(
        &self,
        user_id: i32,
        now: DateTime<Utc>,
    ) -> Result<(), AppError> {
        sqlx::query(
            "UPDATE refresh_tokens
             SET revoked_at = ?
             WHERE user_id = ?
               AND revoked_at IS NULL
               AND expires_at > ?",
        )
        .bind(now)
        .bind(user_id)
        .bind(now)
        .execute(self.0)
        .await
        .map_err(|e| {
            error!("revoke_active_tokens failed (user_id={user_id}): {e}");
            AppError::InternalServerError
        })?;

        Ok(())
    }
}

// ─── テスト ───────────────────────────────────────────────────

/// Repository テスト
///
/// `DATABASE_URL` 環境変数が設定されている場合のみ実行される。
/// CI では `db4test` に向ける。
///
/// # 実行方法
/// ```sh
/// DATABASE_URL="mysql://user:password@127.0.0.1:3306/db4test" cargo test -- --include-ignored
/// ```
#[cfg(test)]
mod tests {
    use chrono::Utc;
    use sqlx::MySqlPool;

    use super::*;

    // ── テスト用プール ────────────────────────────────────────────

    /// `DATABASE_URL` 環境変数からプールを構築する。
    /// 未設定の場合は None を返し、呼び出し元でテストをスキップする。
    async fn test_pool() -> Option<MySqlPool> {
        let url = std::env::var("DATABASE_URL").ok()?;
        MySqlPool::connect(&url).await.ok()
    }

    // ── テスト用ヘルパー ─────────────────────────────────────────

    async fn insert_test_user(pool: &MySqlPool, username: &str, email: &str) -> i32 {
        let now = Utc::now();
        let result = sqlx::query(
            "INSERT INTO users (username, email, password_hash, is_active, last_login_at, created_at, updated_at)
             VALUES (?, ?, 'hash', 1, ?, ?, ?)",
        )
        .bind(username)
        .bind(email)
        .bind(now)
        .bind(now)
        .bind(now)
        .execute(pool)
        .await
        .expect("failed to insert test user");
        result.last_insert_id() as i32
    }

    async fn insert_test_user_points(pool: &MySqlPool, user_id: i32, total_points: i32, level: i32) {
        let now = Utc::now();
        sqlx::query(
            "INSERT INTO user_points (user_id, total_points, level, created_at, updated_at)
             VALUES (?, ?, ?, ?, ?)",
        )
        .bind(user_id)
        .bind(total_points)
        .bind(level)
        .bind(now)
        .bind(now)
        .execute(pool)
        .await
        .expect("failed to insert test user_points");
    }

    async fn insert_test_point_history(
        pool: &MySqlPool,
        user_id: i32,
        points: i32,
        action: &str,
        reference_id: i32,
    ) -> i32 {
        let result = sqlx::query(
            "INSERT INTO point_history (user_id, points, action, reference_id, created_at)
             VALUES (?, ?, ?, ?, NOW())",
        )
        .bind(user_id)
        .bind(points)
        .bind(action)
        .bind(reference_id)
        .execute(pool)
        .await
        .expect("failed to insert test point_history");
        result.last_insert_id() as i32
    }

    /// テストデータの衝突を避けるためユーザー名・メールにサフィックスを付ける
    fn unique(base: &str) -> String {
        format!("{}_{}", base, uuid::Uuid::new_v4().simple())
    }

    // ── UserRepository テスト ─────────────────────────────────────

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_find_by_email_returns_user() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let username = unique("alice");
        let email = unique("alice@example.com");
        let user_id = insert_test_user(&pool, &username, &email).await;
        let repo = MySqlUserRepository(&pool);

        let result = repo.find_by_email(&email).await;

        assert!(result.is_ok());
        let user = result.unwrap().expect("user should exist");
        assert_eq!(user.id, user_id);
        assert_eq!(user.username, username);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_find_by_email_returns_none_for_unknown() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlUserRepository(&pool);

        let result = repo.find_by_email("nobody_unknown@example.com").await;

        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_find_by_id_returns_user() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let username = unique("bob");
        let email = unique("bob@example.com");
        let user_id = insert_test_user(&pool, &username, &email).await;
        let repo = MySqlUserRepository(&pool);

        let result = repo.find_by_id(user_id).await;

        assert!(result.is_ok());
        let user = result.unwrap().expect("user should exist");
        assert_eq!(user.username, username);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_find_by_id_returns_none_for_unknown() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlUserRepository(&pool);

        let result = repo.find_by_id(999999).await;

        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_create_inserts_and_returns_id() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlUserRepository(&pool);
        let username = unique("carol");
        let email = unique("carol@example.com");
        let now = Utc::now();

        let result = repo.create(&username, &email, "hashed_pw", now).await;

        assert!(result.is_ok());
        let user_id = result.unwrap();
        assert!(user_id > 0);
        let found = repo.find_by_id(user_id).await.unwrap().expect("should be found");
        assert_eq!(found.username, username);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_create_returns_conflict_on_duplicate_email() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let username = unique("dave");
        let email = unique("dave@example.com");
        insert_test_user(&pool, &username, &email).await;
        let repo = MySqlUserRepository(&pool);
        let now = Utc::now();

        let result = repo.create(&unique("dave2"), &email, "hashed_pw", now).await;

        assert!(matches!(result, Err(AppError::Conflict(_))));
    }

    // ── PointRepository テスト ────────────────────────────────────

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_by_user_id_returns_row() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool, &unique("eve"), &unique("eve@example.com")).await;
        insert_test_user_points(&pool, user_id, 250, 3).await;
        let repo = MySqlPointRepository(&pool);

        let result = repo.find_by_user_id(user_id).await;

        assert!(result.is_ok());
        let row = result.unwrap().expect("row should exist");
        assert_eq!(row.total_points, 250);
        assert_eq!(row.level, 3);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_by_user_id_returns_none_for_unknown() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlPointRepository(&pool);

        let result = repo.find_by_user_id(999999).await;

        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_history_returns_ordered_records() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool, &unique("frank"), &unique("frank@example.com")).await;
        insert_test_point_history(&pool, user_id, 20, "review", 1).await;
        // タイムスタンプの精度が 1 秒のため、INSERT 間に待機して順序を確定させる
        tokio::time::sleep(std::time::Duration::from_secs(1)).await;
        insert_test_point_history(&pool, user_id, 10, "watch_history", 2).await;
        let repo = MySqlPointRepository(&pool);

        let result = repo.find_history_by_user_id(user_id, 10, 0).await;

        assert!(result.is_ok());
        let rows = result.unwrap();
        assert_eq!(rows.len(), 2);
        // DESC 順なので後に挿入した watch_history が先に来る
        assert_eq!(rows[0].action, "watch_history");
        assert_eq!(rows[1].action, "review");
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_history_respects_limit_and_offset() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool, &unique("grace"), &unique("grace@example.com")).await;
        insert_test_point_history(&pool, user_id, 20, "review", 1).await;
        insert_test_point_history(&pool, user_id, 10, "watch_history", 2).await;
        insert_test_point_history(&pool, user_id, 15, "watch_history", 3).await;
        let repo = MySqlPointRepository(&pool);

        let result = repo.find_history_by_user_id(user_id, 1, 1).await;

        assert!(result.is_ok());
        assert_eq!(result.unwrap().len(), 1);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_history_returns_empty_for_unknown_user() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlPointRepository(&pool);

        let result = repo.find_history_by_user_id(999999, 10, 0).await;

        assert!(result.is_ok());
        assert!(result.unwrap().is_empty());
    }

    // ── TokenRepository テスト ────────────────────────────────────

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_token_save_refresh_token_inserts_record() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool, &unique("henry"), &unique("henry@example.com")).await;
        let repo = MySqlTokenRepository(&pool);
        let now = Utc::now();
        let expires_at = now + chrono::Duration::days(30);
        let token_hash = unique("token_hash");

        let result = repo.save_refresh_token(user_id, &token_hash, expires_at, now).await;

        assert!(result.is_ok());
        let count: i64 = sqlx::query_scalar(
            "SELECT COUNT(*) FROM refresh_tokens WHERE user_id = ? AND token_hash = ?",
        )
        .bind(user_id)
        .bind(&token_hash)
        .fetch_one(&pool)
        .await
        .unwrap();
        assert_eq!(count, 1);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_token_revoke_active_tokens_sets_revoked_at() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool, &unique("ivan"), &unique("ivan@example.com")).await;
        let repo = MySqlTokenRepository(&pool);
        let now = Utc::now();
        let expires_at = now + chrono::Duration::days(30);

        // 有効なトークンを2件保存
        repo.save_refresh_token(user_id, &unique("hash_a"), expires_at, now).await.unwrap();
        repo.save_refresh_token(user_id, &unique("hash_b"), expires_at, now).await.unwrap();

        let result = repo.revoke_active_tokens(user_id, now).await;

        assert!(result.is_ok());
        // revoked_at が設定されたレコードが2件あること
        let revoked_count: i64 = sqlx::query_scalar(
            "SELECT COUNT(*) FROM refresh_tokens WHERE user_id = ? AND revoked_at IS NOT NULL",
        )
        .bind(user_id)
        .fetch_one(&pool)
        .await
        .unwrap();
        assert_eq!(revoked_count, 2);
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_token_revoke_does_not_affect_already_revoked() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let user_id = insert_test_user(&pool, &unique("judy"), &unique("judy@example.com")).await;
        let repo = MySqlTokenRepository(&pool);
        let now = Utc::now();
        let expires_at = now + chrono::Duration::days(30);
        let token_hash = unique("hash_c");

        repo.save_refresh_token(user_id, &token_hash, expires_at, now).await.unwrap();
        // 1回目の revoke
        repo.revoke_active_tokens(user_id, now).await.unwrap();
        // 2回目の revoke は何も変えない（0件更新でエラーにならない）
        let result = repo.revoke_active_tokens(user_id, now).await;

        assert!(result.is_ok());
    }
}
