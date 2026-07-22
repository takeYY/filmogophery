use sqlx::MySqlPool;
use tracing::error;

use crate::app::responses::AppError;

// ─── DB 行型 ──────────────────────────────────────────────────

#[derive(Debug, Clone, sqlx::FromRow)]
pub struct GenreRow {
    pub code: String,
    pub name: String,
}

#[derive(Debug, Clone, sqlx::FromRow)]
pub struct PlatformRow {
    pub id: i32,
    pub code: String,
    pub name: String,
}

// ─── Repository trait ─────────────────────────────────────────

pub trait MasterRepository {
    /// 全ジャンルを取得
    async fn find_all_genres(&self) -> Result<Vec<GenreRow>, AppError>;

    /// 全プラットフォームを取得
    async fn find_all_platforms(&self) -> Result<Vec<PlatformRow>, AppError>;
}

// ─── MySQL 実装 ───────────────────────────────────────────────

pub struct MySqlMasterRepository<'a>(pub &'a MySqlPool);

impl MasterRepository for MySqlMasterRepository<'_> {
    async fn find_all_genres(&self) -> Result<Vec<GenreRow>, AppError> {
        sqlx::query_as("SELECT code, name FROM genres ORDER BY code")
            .fetch_all(self.0)
            .await
            .map_err(|e| {
                error!("find_all_genres failed: {e}");
                AppError::InternalServerError
            })
    }

    async fn find_all_platforms(&self) -> Result<Vec<PlatformRow>, AppError> {
        sqlx::query_as("SELECT id, code, name FROM platforms ORDER BY id")
            .fetch_all(self.0)
            .await
            .map_err(|e| {
                error!("find_all_platforms failed: {e}");
                AppError::InternalServerError
            })
    }
}

// ─── テスト ───────────────────────────────────────────────────

/// Repository テスト
///
/// `DATABASE_URL` 環境変数が設定されている場合のみ実行される。
///
/// # 実行方法
/// ```sh
/// DATABASE_URL="mysql://user:password@127.0.0.1:3306/db4test" cargo test -- --include-ignored
/// ```
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
    async fn test_find_all_genres_returns_list() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMasterRepository(&pool);

        let result = repo.find_all_genres().await;

        assert!(result.is_ok());
        // genres テーブルにレコードが存在すること（シードデータ前提）
        assert!(!result.unwrap().is_empty());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_find_all_platforms_returns_list() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMasterRepository(&pool);

        let result = repo.find_all_platforms().await;

        assert!(result.is_ok());
        // platforms テーブルにレコードが存在すること（シードデータ前提）
        assert!(!result.unwrap().is_empty());
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_find_all_genres_rows_have_code_and_name() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMasterRepository(&pool);

        let genres = repo.find_all_genres().await.unwrap();

        for genre in &genres {
            assert!(!genre.code.is_empty(), "code should not be empty");
            assert!(!genre.name.is_empty(), "name should not be empty");
        }
    }

    #[tokio::test]
    #[ignore = "requires DATABASE_URL"]
    async fn test_find_all_platforms_rows_have_id_code_and_name() {
        let pool = test_pool().await.expect("DATABASE_URL not set");
        let repo = MySqlMasterRepository(&pool);

        let platforms = repo.find_all_platforms().await.unwrap();

        for platform in &platforms {
            assert!(platform.id > 0, "id should be positive");
            assert!(!platform.code.is_empty(), "code should not be empty");
            assert!(!platform.name.is_empty(), "name should not be empty");
        }
    }
}
