use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use validator::Validate;

use crate::app::responses::{AppError, ApiResult};

use super::repository::{PointRepository, TokenRepository, UserRepository};

// ─── レスポンス型 ─────────────────────────────────────────────

/// POST /v1/users レスポンス
#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TokenResponse {
    pub access_token: String,
    pub refresh_token: String,
    pub token_type: String,
    pub expires_in: i64,
    pub expires_at: DateTime<Utc>,
}

/// GET /v1/users/me レスポンス
#[derive(Debug, Serialize)]
pub struct UserResponse {
    pub id: i32,
    pub username: String,
}

/// GET /v1/users/me/points レスポンス
#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct UserPointsResponse {
    pub total_points: i32,
    pub level: i32,
    pub next_level_points: i32,
    pub current_level_width: i32,
    pub point_history: Vec<PointHistoryItem>,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct PointHistoryItem {
    pub id: i32,
    pub points: i32,
    pub action: String,
    pub reference_id: i32,
    pub created_at: Option<DateTime<Utc>>,
}

// ─── リクエスト型 ─────────────────────────────────────────────

#[derive(Debug, Deserialize, Validate)]
pub struct CreateUserInput {
    #[validate(length(min = 1, message = "username is required"))]
    pub username: String,
    #[validate(email(message = "email is invalid"))]
    pub email: String,
    #[validate(length(min = 8, message = "password must be at least 8 characters"))]
    pub password: String,
}

// ─── レベル計算ロジック ────────────────────────────────────────

/// レベルの閾値（累計ポイント）
/// Lv1→2: 100pt, Lv2→3: 200pt, Lv3→4: 400pt, Lv4→5: 800pt
const LEVEL_THRESHOLDS: [i32; 5] = [0, 100, 300, 700, 1500];
const FIXED_LEVEL_POINTS: i32 = 1000;

/// 累計ポイントからレベルを計算する（レビュー・視聴履歴のポイント付与時にも使用）
#[allow(dead_code)]
pub fn calc_level(total_points: i32) -> i32 {
    for i in (0..LEVEL_THRESHOLDS.len()).rev() {
        if total_points >= LEVEL_THRESHOLDS[i] {
            let base_level = (i + 1) as i32;
            if base_level < LEVEL_THRESHOLDS.len() as i32 {
                return base_level;
            }
            let extra =
                (total_points - LEVEL_THRESHOLDS[LEVEL_THRESHOLDS.len() - 1]) / FIXED_LEVEL_POINTS;
            return LEVEL_THRESHOLDS.len() as i32 + extra;
        }
    }
    1
}

pub fn calc_next_level_points(total_points: i32, level: i32) -> i32 {
    if (level as usize) < LEVEL_THRESHOLDS.len() {
        return LEVEL_THRESHOLDS[level as usize] - total_points;
    }
    let points_in_current =
        (total_points - LEVEL_THRESHOLDS[LEVEL_THRESHOLDS.len() - 1]) % FIXED_LEVEL_POINTS;
    FIXED_LEVEL_POINTS - points_in_current
}

pub fn calc_current_level_width(level: i32) -> i32 {
    if (level as usize) < LEVEL_THRESHOLDS.len() {
        return LEVEL_THRESHOLDS[level as usize] - LEVEL_THRESHOLDS[(level - 1) as usize];
    }
    FIXED_LEVEL_POINTS
}

// ─── Use Case 関数 ─────────────────────────────────────────────

/// POST /v1/users — ユーザー登録
pub async fn create_user<UR, TR>(
    user_repo: &UR,
    token_repo: &TR,
    jwt_secret: &str,
    input: CreateUserInput,
) -> ApiResult<TokenResponse>
where
    UR: UserRepository,
    TR: TokenRepository,
{
    // メールの重複チェック
    if user_repo.find_by_email(&input.email).await?.is_some() {
        return Err(AppError::Conflict("user".to_string()));
    }

    // パスワードをハッシュ化
    let pwd_hash = bcrypt::hash(&input.password, bcrypt::DEFAULT_COST)
        .map_err(|_| AppError::InternalServerError)?;

    let now = Utc::now();
    let user_id = user_repo.create(&input.username, &input.email, &pwd_hash, now).await?;

    issue_token(token_repo, jwt_secret, user_id, now).await
}

/// GET /v1/users/me — ログインユーザー取得
pub async fn get_current_user<UR>(user_repo: &UR, user_id: i32) -> ApiResult<UserResponse>
where
    UR: UserRepository,
{
    let user = user_repo
        .find_by_id(user_id)
        .await?
        .ok_or_else(|| AppError::NotFound("user".to_string()))?;

    Ok(UserResponse {
        id: user.id,
        username: user.username,
    })
}

/// GET /v1/users/me/points — ポイント・レベル取得
pub async fn get_user_points<PR>(
    point_repo: &PR,
    user_id: i32,
    limit: i32,
    offset: i32,
) -> ApiResult<UserPointsResponse>
where
    PR: PointRepository,
{
    let (total_points, level) = match point_repo.find_by_user_id(user_id).await? {
        Some(row) => (row.total_points, row.level),
        None => (0, 1),
    };

    let histories = point_repo
        .find_history_by_user_id(user_id, limit, offset)
        .await?;

    let next_level_points = calc_next_level_points(total_points, level);
    let current_level_width = calc_current_level_width(level);

    let point_history = histories
        .into_iter()
        .map(|h| PointHistoryItem {
            id: h.id,
            points: h.points,
            action: h.action,
            reference_id: h.reference_id,
            created_at: h.created_at,
        })
        .collect();

    Ok(UserPointsResponse {
        total_points,
        level,
        next_level_points,
        current_level_width,
        point_history,
    })
}

// ─── 内部ヘルパー ─────────────────────────────────────────────

/// JWT + リフレッシュトークンを発行して TokenResponse を返す
async fn issue_token<TR>(
    token_repo: &TR,
    jwt_secret: &str,
    user_id: i32,
    now: DateTime<Utc>,
) -> ApiResult<TokenResponse>
where
    TR: TokenRepository,
{
    use crate::pkg::jwt;

    let expires_in: i64 = 3600;
    let access_token = jwt::generate_access_token(user_id, expires_in, jwt_secret)?;

    let refresh_token = uuid::Uuid::new_v4().to_string();
    let refresh_token_hash = bcrypt::hash(&refresh_token, bcrypt::DEFAULT_COST)
        .map_err(|_| AppError::InternalServerError)?;

    let refresh_expires_at = now + chrono::Duration::days(30);
    token_repo
        .save_refresh_token(user_id, &refresh_token_hash, refresh_expires_at, now)
        .await?;

    Ok(TokenResponse {
        access_token,
        refresh_token,
        token_type: "Bearer".to_string(),
        expires_in,
        expires_at: now + chrono::Duration::seconds(expires_in),
    })
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use std::sync::Mutex;

    use super::*;
    use crate::app::features::user::repository::{
        PointHistoryRow, PointRepository, TokenRepository, UserRepository, UserRow,
        UserPointsRow,
    };

    // ── モック: UserRepository ──────────────────────────────────

    struct MockUserRepository {
        /// find_by_email / find_by_id で返すユーザー（None でユーザーなし）
        user: Option<UserRow>,
        /// create で返す user_id
        created_id: i32,
    }

    impl MockUserRepository {
        fn with_user(id: i32, username: &str) -> Self {
            Self {
                user: Some(UserRow {
                    id,
                    username: username.to_string(),
                    password_hash: "hash".to_string(),
                }),
                created_id: id,
            }
        }

        fn empty(next_id: i32) -> Self {
            Self {
                user: None,
                created_id: next_id,
            }
        }
    }

    impl UserRepository for MockUserRepository {
        async fn find_by_email(&self, _email: &str) -> Result<Option<UserRow>, AppError> {
            Ok(self.user.clone())
        }

        async fn find_by_id(&self, _id: i32) -> Result<Option<UserRow>, AppError> {
            Ok(self.user.clone())
        }

        async fn create(
            &self,
            _username: &str,
            _email: &str,
            _password_hash: &str,
            _now: DateTime<Utc>,
        ) -> Result<i32, AppError> {
            Ok(self.created_id)
        }
    }

    // UserRow に Clone が必要なのでここで派生させる
    #[derive(Clone)]
    struct ClonableUserRow {
        id: i32,
        username: String,
        password_hash: String,
    }

    impl From<ClonableUserRow> for UserRow {
        fn from(r: ClonableUserRow) -> Self {
            UserRow {
                id: r.id,
                username: r.username,
                password_hash: r.password_hash,
            }
        }
    }

    // UserRow 自体に Clone を付けるため repository.rs を変えるより
    // モック内で Option<UserRow> を都度構築する形にする
    // ── モック: PointRepository ─────────────────────────────────

    struct MockPointRepository {
        points_row: Option<UserPointsRow>,
        history: Vec<PointHistoryRow>,
    }

    impl MockPointRepository {
        fn with_points(total: i32, level: i32, history: Vec<PointHistoryRow>) -> Self {
            Self {
                points_row: Some(UserPointsRow {
                    total_points: total,
                    level,
                }),
                history,
            }
        }

        fn empty() -> Self {
            Self {
                points_row: None,
                history: vec![],
            }
        }
    }

    impl PointRepository for MockPointRepository {
        async fn find_by_user_id(
            &self,
            _user_id: i32,
        ) -> Result<Option<UserPointsRow>, AppError> {
            Ok(self.points_row.as_ref().map(|r| UserPointsRow {
                total_points: r.total_points,
                level: r.level,
            }))
        }

        async fn find_history_by_user_id(
            &self,
            _user_id: i32,
            _limit: i32,
            _offset: i32,
        ) -> Result<Vec<PointHistoryRow>, AppError> {
            Ok(self.history.iter().map(|h| PointHistoryRow {
                id: h.id,
                points: h.points,
                action: h.action.clone(),
                reference_id: h.reference_id,
                created_at: h.created_at,
            }).collect())
        }
    }

    // ── モック: TokenRepository ─────────────────────────────────

    struct MockTokenRepository {
        /// save_refresh_token が呼ばれた回数を記録
        call_count: Mutex<u32>,
    }

    impl MockTokenRepository {
        fn new() -> Self {
            Self {
                call_count: Mutex::new(0),
            }
        }

        fn call_count(&self) -> u32 {
            *self.call_count.lock().unwrap()
        }
    }

    impl TokenRepository for MockTokenRepository {
        async fn save_refresh_token(
            &self,
            _user_id: i32,
            _token_hash: &str,
            _expires_at: DateTime<Utc>,
            _now: DateTime<Utc>,
        ) -> Result<(), AppError> {
            *self.call_count.lock().unwrap() += 1;
            Ok(())
        }

        async fn revoke_active_tokens(
            &self,
            _user_id: i32,
            _now: DateTime<Utc>,
        ) -> Result<(), AppError> {
            Ok(())
        }
    }

    // ── レベル計算テスト ─────────────────────────────────────────

    #[test]
    fn test_calc_level_initial() {
        assert_eq!(calc_level(0), 1);
    }

    #[test]
    fn test_calc_level_boundaries() {
        // Lv1→2: 100pt
        assert_eq!(calc_level(99), 1);
        assert_eq!(calc_level(100), 2);
        // Lv2→3: 300pt
        assert_eq!(calc_level(299), 2);
        assert_eq!(calc_level(300), 3);
        // Lv3→4: 700pt
        assert_eq!(calc_level(699), 3);
        assert_eq!(calc_level(700), 4);
        // Lv4→5: 1500pt
        assert_eq!(calc_level(1499), 4);
        assert_eq!(calc_level(1500), 5);
    }

    #[test]
    fn test_calc_level_high_points() {
        // Lv5以降: 1000pt刻み
        assert_eq!(calc_level(2499), 5);
        assert_eq!(calc_level(2500), 6);
        assert_eq!(calc_level(3499), 6);
        assert_eq!(calc_level(3500), 7);
    }

    #[test]
    fn test_calc_next_level_points() {
        // Lv1(0pt): あと100pt
        assert_eq!(calc_next_level_points(0, 1), 100);
        // Lv1(50pt): あと50pt
        assert_eq!(calc_next_level_points(50, 1), 50);
        // Lv2(100pt): あと200pt
        assert_eq!(calc_next_level_points(100, 2), 200);
        // Lv5(1500pt): 次の1000pt区切りまで1000pt
        assert_eq!(calc_next_level_points(1500, 5), 1000);
        // Lv5(1700pt): あと800pt
        assert_eq!(calc_next_level_points(1700, 5), 800);
    }

    #[test]
    fn test_calc_current_level_width() {
        assert_eq!(calc_current_level_width(1), 100);  // Lv1幅
        assert_eq!(calc_current_level_width(2), 200);  // Lv2幅
        assert_eq!(calc_current_level_width(3), 400);  // Lv3幅
        assert_eq!(calc_current_level_width(4), 800);  // Lv4幅
        assert_eq!(calc_current_level_width(5), 1000); // Lv5以降は固定
        assert_eq!(calc_current_level_width(6), 1000);
    }

    // ── create_user テスト ────────────────────────────────────────

    #[tokio::test]
    async fn test_create_user_success() {
        let user_repo = MockUserRepository::empty(1);
        let token_repo = MockTokenRepository::new();

        let input = CreateUserInput {
            username: "testuser".to_string(),
            email: "test@example.com".to_string(),
            password: "password123".to_string(),
        };

        let result = create_user(&user_repo, &token_repo, "secret", input).await;

        assert!(result.is_ok());
        let token = result.unwrap();
        assert_eq!(token.token_type, "Bearer");
        assert_eq!(token.expires_in, 3600);
        assert!(!token.access_token.is_empty());
        assert!(!token.refresh_token.is_empty());
        // refresh token が保存されていること
        assert_eq!(token_repo.call_count(), 1);
    }

    #[tokio::test]
    async fn test_create_user_conflict_when_email_exists() {
        // メール重複: find_by_email がユーザーを返す
        let user_repo = MockUserRepository::with_user(1, "existing");
        let token_repo = MockTokenRepository::new();

        let input = CreateUserInput {
            username: "newuser".to_string(),
            email: "existing@example.com".to_string(),
            password: "password123".to_string(),
        };

        let result = create_user(&user_repo, &token_repo, "secret", input).await;

        assert!(matches!(result, Err(AppError::Conflict(_))));
        // 重複時は refresh token を保存しないこと
        assert_eq!(token_repo.call_count(), 0);
    }

    // ── get_current_user テスト ───────────────────────────────────

    #[tokio::test]
    async fn test_get_current_user_success() {
        let user_repo = MockUserRepository::with_user(42, "alice");

        let result = get_current_user(&user_repo, 42).await;

        assert!(result.is_ok());
        let user = result.unwrap();
        assert_eq!(user.id, 42);
        assert_eq!(user.username, "alice");
    }

    #[tokio::test]
    async fn test_get_current_user_not_found() {
        let user_repo = MockUserRepository::empty(0);

        let result = get_current_user(&user_repo, 99).await;

        assert!(matches!(result, Err(AppError::NotFound(_))));
    }

    // ── get_user_points テスト ────────────────────────────────────

    #[tokio::test]
    async fn test_get_user_points_with_existing_record() {
        let point_repo = MockPointRepository::with_points(150, 2, vec![
            PointHistoryRow {
                id: 1,
                points: 20,
                action: "review".to_string(),
                reference_id: 10,
                created_at: None,
            },
        ]);

        let result = get_user_points(&point_repo, 1, 20, 0).await;

        assert!(result.is_ok());
        let resp = result.unwrap();
        assert_eq!(resp.total_points, 150);
        assert_eq!(resp.level, 2);
        // Lv2(150pt): 300-150=150pt
        assert_eq!(resp.next_level_points, 150);
        // Lv2幅: 200pt
        assert_eq!(resp.current_level_width, 200);
        assert_eq!(resp.point_history.len(), 1);
        assert_eq!(resp.point_history[0].action, "review");
    }

    #[tokio::test]
    async fn test_get_user_points_no_record_returns_defaults() {
        // user_points レコードがない場合はデフォルト(0pt, Lv1)
        let point_repo = MockPointRepository::empty();

        let result = get_user_points(&point_repo, 1, 20, 0).await;

        assert!(result.is_ok());
        let resp = result.unwrap();
        assert_eq!(resp.total_points, 0);
        assert_eq!(resp.level, 1);
        assert_eq!(resp.next_level_points, 100);
        assert_eq!(resp.current_level_width, 100);
        assert!(resp.point_history.is_empty());
    }

    #[tokio::test]
    async fn test_get_user_points_history_is_empty_when_no_history() {
        let point_repo = MockPointRepository::with_points(50, 1, vec![]);

        let result = get_user_points(&point_repo, 1, 20, 0).await;

        assert!(result.is_ok());
        assert!(result.unwrap().point_history.is_empty());
    }
}
