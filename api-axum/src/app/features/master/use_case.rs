use serde::Serialize;

use crate::app::responses::{ApiResult, AppError};

use super::repository::MasterRepository;

// ─── レスポンス型 ─────────────────────────────────────────────

/// GET /v1/genres レスポンスアイテム
#[derive(Debug, Serialize)]
pub struct GenreResponse {
    pub code: String,
    pub name: String,
}

/// GET /v1/platforms レスポンスアイテム
#[derive(Debug, Serialize)]
pub struct PlatformResponse {
    pub id: i32,
    pub code: String,
    pub name: String,
}

// ─── Use Case 関数 ─────────────────────────────────────────────

/// GET /v1/genres — ジャンル一覧取得
pub async fn get_genres<R>(repo: &R) -> ApiResult<Vec<GenreResponse>>
where
    R: MasterRepository,
{
    let rows = repo.find_all_genres().await?;
    Ok(rows
        .into_iter()
        .map(|r| GenreResponse {
            code: r.code,
            name: r.name,
        })
        .collect())
}

/// GET /v1/platforms — プラットフォーム一覧取得
pub async fn get_platforms<R>(repo: &R) -> ApiResult<Vec<PlatformResponse>>
where
    R: MasterRepository,
{
    let rows = repo.find_all_platforms().await?;
    Ok(rows
        .into_iter()
        .map(|r| PlatformResponse {
            id: r.id,
            code: r.code,
            name: r.name,
        })
        .collect())
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use super::*;
    use crate::app::features::master::repository::{GenreRow, MasterRepository, PlatformRow};

    // ── モック ────────────────────────────────────────────────────

    struct MockMasterRepository {
        genres: Result<Vec<GenreRow>, AppError>,
        platforms: Result<Vec<PlatformRow>, AppError>,
    }

    impl MockMasterRepository {
        fn with_data(genres: Vec<GenreRow>, platforms: Vec<PlatformRow>) -> Self {
            Self {
                genres: Ok(genres),
                platforms: Ok(platforms),
            }
        }

        fn empty() -> Self {
            Self {
                genres: Ok(vec![]),
                platforms: Ok(vec![]),
            }
        }

        fn failing() -> Self {
            Self {
                genres: Err(AppError::InternalServerError),
                platforms: Err(AppError::InternalServerError),
            }
        }
    }

    impl MasterRepository for MockMasterRepository {
        async fn find_all_genres(&self) -> Result<Vec<GenreRow>, AppError> {
            match &self.genres {
                Ok(rows) => Ok(rows.iter().map(|r| GenreRow {
                    code: r.code.clone(),
                    name: r.name.clone(),
                }).collect()),
                Err(_) => Err(AppError::InternalServerError),
            }
        }

        async fn find_all_platforms(&self) -> Result<Vec<PlatformRow>, AppError> {
            match &self.platforms {
                Ok(rows) => Ok(rows.iter().map(|r| PlatformRow {
                    id: r.id,
                    code: r.code.clone(),
                    name: r.name.clone(),
                }).collect()),
                Err(_) => Err(AppError::InternalServerError),
            }
        }
    }

    // ── get_genres テスト ─────────────────────────────────────────

    #[tokio::test]
    async fn test_get_genres_returns_mapped_response() {
        let repo = MockMasterRepository::with_data(
            vec![
                GenreRow { code: "ACTION".to_string(), name: "アクション".to_string() },
                GenreRow { code: "DRAMA".to_string(),  name: "ドラマ".to_string() },
            ],
            vec![],
        );

        let result = get_genres(&repo).await.unwrap();

        assert_eq!(result.len(), 2);
        assert_eq!(result[0].code, "ACTION");
        assert_eq!(result[0].name, "アクション");
        assert_eq!(result[1].code, "DRAMA");
        assert_eq!(result[1].name, "ドラマ");
    }

    #[tokio::test]
    async fn test_get_genres_returns_empty_when_no_data() {
        let repo = MockMasterRepository::empty();

        let result = get_genres(&repo).await.unwrap();

        assert!(result.is_empty());
    }

    #[tokio::test]
    async fn test_get_genres_propagates_repo_error() {
        let repo = MockMasterRepository::failing();

        let result = get_genres(&repo).await;

        assert!(matches!(result, Err(AppError::InternalServerError)));
    }

    // ── get_platforms テスト ──────────────────────────────────────

    #[tokio::test]
    async fn test_get_platforms_returns_mapped_response() {
        let repo = MockMasterRepository::with_data(
            vec![],
            vec![
                PlatformRow { id: 1, code: "NETFLIX".to_string(), name: "Netflix".to_string() },
                PlatformRow { id: 2, code: "PRIME".to_string(),   name: "Prime Video".to_string() },
            ],
        );

        let result = get_platforms(&repo).await.unwrap();

        assert_eq!(result.len(), 2);
        assert_eq!(result[0].id, 1);
        assert_eq!(result[0].code, "NETFLIX");
        assert_eq!(result[0].name, "Netflix");
        assert_eq!(result[1].id, 2);
        assert_eq!(result[1].code, "PRIME");
    }

    #[tokio::test]
    async fn test_get_platforms_returns_empty_when_no_data() {
        let repo = MockMasterRepository::empty();

        let result = get_platforms(&repo).await.unwrap();

        assert!(result.is_empty());
    }

    #[tokio::test]
    async fn test_get_platforms_propagates_repo_error() {
        let repo = MockMasterRepository::failing();

        let result = get_platforms(&repo).await;

        assert!(matches!(result, Err(AppError::InternalServerError)));
    }
}
