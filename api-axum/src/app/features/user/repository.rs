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
}

// ─── テスト ───────────────────────────────────────────────────

/// Repository テスト
///
/// `sqlx::test` マクロが自動的に `TEST_DATABASE_URL` 環境変数に接続した
/// 一時データベースを用意し、テスト終了時にロールバックする。
/// CI では `TEST_DATABASE_URL` を MySQL の db4test に向ける。
///
/// # 実行方法
/// ```sh
/// TEST_DATABASE_URL="mysql://user:password@127.0.0.1:3306/db4test" cargo test
/// ```
#[cfg(test)]
mod tests {
    use chrono::Utc;
    use sqlx::MySqlPool;

    use super::*;

    // ── テスト用ヘルパー ─────────────────────────────────────────

    /// テスト用ユーザーを INSERT して id を返す
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

    /// テスト用 user_points レコードを INSERT する
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

    /// テスト用 point_history レコードを INSERT する
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

    // ── UserRepository テスト ─────────────────────────────────────

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_find_by_email_returns_user(pool: MySqlPool) {
        let user_id = insert_test_user(&pool, "alice", "alice@example.com").await;
        let repo = MySqlUserRepository(&pool);

        let result = repo.find_by_email("alice@example.com").await;

        assert!(result.is_ok());
        let user = result.unwrap();
        assert!(user.is_some());
        let user = user.unwrap();
        assert_eq!(user.id, user_id);
        assert_eq!(user.username, "alice");
    }

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_find_by_email_returns_none_for_unknown(pool: MySqlPool) {
        let repo = MySqlUserRepository(&pool);

        let result = repo.find_by_email("nobody@example.com").await;

        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_find_by_id_returns_user(pool: MySqlPool) {
        let user_id = insert_test_user(&pool, "bob", "bob@example.com").await;
        let repo = MySqlUserRepository(&pool);

        let result = repo.find_by_id(user_id).await;

        assert!(result.is_ok());
        let user = result.unwrap();
        assert!(user.is_some());
        assert_eq!(user.unwrap().username, "bob");
    }

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_find_by_id_returns_none_for_unknown(pool: MySqlPool) {
        let repo = MySqlUserRepository(&pool);

        let result = repo.find_by_id(999999).await;

        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_create_inserts_and_returns_id(pool: MySqlPool) {
        let repo = MySqlUserRepository(&pool);
        let now = Utc::now();

        let result = repo
            .create("carol", "carol@example.com", "hashed_pw", now)
            .await;

        assert!(result.is_ok());
        let user_id = result.unwrap();
        assert!(user_id > 0);

        // 実際に挿入されていること
        let found = repo.find_by_id(user_id).await.unwrap();
        assert!(found.is_some());
        assert_eq!(found.unwrap().username, "carol");
    }

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_user_create_returns_conflict_on_duplicate_email(pool: MySqlPool) {
        insert_test_user(&pool, "dave", "dave@example.com").await;
        let repo = MySqlUserRepository(&pool);
        let now = Utc::now();

        // 同じメールで再登録
        let result = repo
            .create("dave2", "dave@example.com", "hashed_pw", now)
            .await;

        assert!(matches!(result, Err(AppError::Conflict(_))));
    }

    // ── PointRepository テスト ────────────────────────────────────

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_by_user_id_returns_row(pool: MySqlPool) {
        let user_id = insert_test_user(&pool, "eve", "eve@example.com").await;
        insert_test_user_points(&pool, user_id, 250, 3).await;
        let repo = MySqlPointRepository(&pool);

        let result = repo.find_by_user_id(user_id).await;

        assert!(result.is_ok());
        let row = result.unwrap();
        assert!(row.is_some());
        let row = row.unwrap();
        assert_eq!(row.total_points, 250);
        assert_eq!(row.level, 3);
    }

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_by_user_id_returns_none_for_unknown(pool: MySqlPool) {
        let repo = MySqlPointRepository(&pool);

        let result = repo.find_by_user_id(999999).await;

        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_history_returns_ordered_records(pool: MySqlPool) {
        let user_id = insert_test_user(&pool, "frank", "frank@example.com").await;
        // 2件挿入（created_at の順序は INSERT 順に依存）
        insert_test_point_history(&pool, user_id, 20, "review", 1).await;
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

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_history_respects_limit_and_offset(pool: MySqlPool) {
        let user_id = insert_test_user(&pool, "grace", "grace@example.com").await;
        insert_test_point_history(&pool, user_id, 20, "review", 1).await;
        insert_test_point_history(&pool, user_id, 10, "watch_history", 2).await;
        insert_test_point_history(&pool, user_id, 15, "watch_history", 3).await;
        let repo = MySqlPointRepository(&pool);

        // limit=1, offset=1 → 2番目の1件だけ
        let result = repo.find_history_by_user_id(user_id, 1, 1).await;

        assert!(result.is_ok());
        assert_eq!(result.unwrap().len(), 1);
    }

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_point_find_history_returns_empty_for_unknown_user(pool: MySqlPool) {
        let repo = MySqlPointRepository(&pool);

        let result = repo.find_history_by_user_id(999999, 10, 0).await;

        assert!(result.is_ok());
        assert!(result.unwrap().is_empty());
    }

    // ── TokenRepository テスト ────────────────────────────────────

    #[sqlx::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_token_save_refresh_token_inserts_record(pool: MySqlPool) {
        let user_id = insert_test_user(&pool, "henry", "henry@example.com").await;
        let repo = MySqlTokenRepository(&pool);
        let now = Utc::now();
        let expires_at = now + chrono::Duration::days(30);

        let result = repo
            .save_refresh_token(user_id, "token_hash_value", expires_at, now)
            .await;

        assert!(result.is_ok());

        // DB に実際に保存されていること
        let count: i64 = sqlx::query_scalar(
            "SELECT COUNT(*) FROM refresh_tokens WHERE user_id = ? AND token_hash = ?",
        )
        .bind(user_id)
        .bind("token_hash_value")
        .fetch_one(&pool)
        .await
        .unwrap();
        assert_eq!(count, 1);
    }
}
