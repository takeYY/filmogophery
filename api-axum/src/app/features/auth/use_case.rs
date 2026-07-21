use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use validator::Validate;

use crate::app::features::user::repository::{TokenRepository, UserRepository};
use crate::app::responses::{ApiResult, AppError};

// ─── レスポンス型 ─────────────────────────────────────────────

/// POST /v1/auth/login レスポンス
#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct TokenResponse {
    pub access_token: String,
    pub refresh_token: String,
    pub token_type: String,
    pub expires_in: i64,
    pub expires_at: DateTime<Utc>,
}

// ─── リクエスト型 ─────────────────────────────────────────────

#[derive(Debug, Deserialize, Validate)]
pub struct LoginInput {
    #[validate(email(message = "email is invalid"))]
    pub email: String,
    #[validate(length(min = 8, message = "password must be at least 8 characters"))]
    pub password: String,
}

// ─── Use Case 関数 ────────────────────────────────────────────

/// POST /v1/auth/login — ログイン
pub async fn login<UR, TR>(
    user_repo: &UR,
    token_repo: &TR,
    jwt_secret: &str,
    input: LoginInput,
) -> ApiResult<TokenResponse>
where
    UR: UserRepository,
    TR: TokenRepository,
{
    // メールアドレスでユーザーを取得
    let user = user_repo
        .find_by_email(&input.email)
        .await?
        .ok_or(AppError::Unauthorized)?;

    // パスワードを検証
    let is_valid = bcrypt::verify(&input.password, &user.password_hash)
        .map_err(|_| AppError::InternalServerError)?;
    if !is_valid {
        return Err(AppError::Unauthorized);
    }

    let now = Utc::now();
    issue_token(token_repo, jwt_secret, user.id, now).await
}

/// POST /v1/auth/logout — ログアウト
pub async fn logout<TR>(
    token_repo: &TR,
    user_id: i32,
) -> ApiResult<()>
where
    TR: TokenRepository,
{
    let now = Utc::now();
    token_repo.revoke_active_tokens(user_id, now).await?;
    Ok(())
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
    use sha2::{Digest, Sha256};

    let expires_in: i64 = 3600;
    let access_token = crate::pkg::jwt::generate_access_token(user_id, expires_in, jwt_secret)?;

    // リフレッシュトークン（UUID）を生成し、SHA-256 ハッシュで DB に保存
    let refresh_token = uuid::Uuid::new_v4().to_string();
    let token_hash = hex::encode(Sha256::digest(refresh_token.as_bytes()));

    let refresh_expires_at = now + chrono::Duration::days(30);
    token_repo
        .save_refresh_token(user_id, &token_hash, refresh_expires_at, now)
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

    use chrono::{DateTime, Utc};

    use super::*;
    use crate::app::features::user::repository::{TokenRepository, UserRepository, UserRow};

    // ── モック: UserRepository ──────────────────────────────────

    struct MockUserRepository {
        user: Option<UserRow>,
    }

    impl MockUserRepository {
        fn with_user(id: i32, password_hash: &str) -> Self {
            Self {
                user: Some(UserRow {
                    id,
                    username: "testuser".to_string(),
                    password_hash: password_hash.to_string(),
                }),
            }
        }

        fn empty() -> Self {
            Self { user: None }
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
            Ok(1)
        }
    }

    // ── モック: TokenRepository ─────────────────────────────────

    struct MockTokenRepository {
        save_call_count: Mutex<u32>,
        revoke_call_count: Mutex<u32>,
        /// save_refresh_token を失敗させるか
        save_fails: bool,
        /// revoke_active_tokens を失敗させるか
        revoke_fails: bool,
    }

    impl MockTokenRepository {
        fn new() -> Self {
            Self {
                save_call_count: Mutex::new(0),
                revoke_call_count: Mutex::new(0),
                save_fails: false,
                revoke_fails: false,
            }
        }

        fn failing_save() -> Self {
            Self { save_fails: true, ..Self::new() }
        }

        fn failing_revoke() -> Self {
            Self { revoke_fails: true, ..Self::new() }
        }

        fn save_count(&self) -> u32 {
            *self.save_call_count.lock().unwrap()
        }

        fn revoke_count(&self) -> u32 {
            *self.revoke_call_count.lock().unwrap()
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
            if self.save_fails {
                return Err(AppError::InternalServerError);
            }
            *self.save_call_count.lock().unwrap() += 1;
            Ok(())
        }

        async fn revoke_active_tokens(
            &self,
            _user_id: i32,
            _now: DateTime<Utc>,
        ) -> Result<(), AppError> {
            if self.revoke_fails {
                return Err(AppError::InternalServerError);
            }
            *self.revoke_call_count.lock().unwrap() += 1;
            Ok(())
        }
    }

    // ── login テスト ─────────────────────────────────────────────

    #[tokio::test]
    async fn test_login_success() {
        // bcrypt でパスワードハッシュを生成
        let password = "password123";
        let hash = bcrypt::hash(password, bcrypt::DEFAULT_COST).unwrap();

        let user_repo = MockUserRepository::with_user(1, &hash);
        let token_repo = MockTokenRepository::new();

        let input = LoginInput {
            email: "test@example.com".to_string(),
            password: password.to_string(),
        };

        let result = login(&user_repo, &token_repo, "secret", input).await;

        assert!(result.is_ok());
        let token = result.unwrap();
        assert_eq!(token.token_type, "Bearer");
        assert_eq!(token.expires_in, 3600);
        assert!(!token.access_token.is_empty());
        assert!(!token.refresh_token.is_empty());
        assert_eq!(token_repo.save_count(), 1);
    }

    #[tokio::test]
    async fn test_login_fails_when_user_not_found() {
        let user_repo = MockUserRepository::empty();
        let token_repo = MockTokenRepository::new();

        let input = LoginInput {
            email: "nobody@example.com".to_string(),
            password: "password123".to_string(),
        };

        let result = login(&user_repo, &token_repo, "secret", input).await;

        assert!(matches!(result, Err(AppError::Unauthorized)));
        assert_eq!(token_repo.save_count(), 0);
    }

    #[tokio::test]
    async fn test_login_fails_when_password_wrong() {
        let hash = bcrypt::hash("correct_password", bcrypt::DEFAULT_COST).unwrap();
        let user_repo = MockUserRepository::with_user(1, &hash);
        let token_repo = MockTokenRepository::new();

        let input = LoginInput {
            email: "test@example.com".to_string(),
            password: "wrong_password".to_string(),
        };

        let result = login(&user_repo, &token_repo, "secret", input).await;

        assert!(matches!(result, Err(AppError::Unauthorized)));
        assert_eq!(token_repo.save_count(), 0);
    }

    // ── logout テスト ────────────────────────────────────────────

    #[tokio::test]
    async fn test_logout_revokes_tokens() {
        let token_repo = MockTokenRepository::new();

        let result = logout(&token_repo, 42).await;

        assert!(result.is_ok());
        assert_eq!(token_repo.revoke_count(), 1);
    }

    #[tokio::test]
    async fn test_logout_propagates_repo_error() {
        let token_repo = MockTokenRepository::failing_revoke();

        let result = logout(&token_repo, 42).await;

        assert!(matches!(result, Err(AppError::InternalServerError)));
    }

    // ── issue_token (login 経由) ──────────────────────────────────

    #[tokio::test]
    async fn test_login_propagates_save_token_error() {
        let password = "password123";
        let hash = bcrypt::hash(password, bcrypt::DEFAULT_COST).unwrap();
        let user_repo = MockUserRepository::with_user(1, &hash);
        let token_repo = MockTokenRepository::failing_save();

        let input = LoginInput {
            email: "test@example.com".to_string(),
            password: password.to_string(),
        };

        let result = login(&user_repo, &token_repo, "secret", input).await;

        assert!(matches!(result, Err(AppError::InternalServerError)));
    }

    #[tokio::test]
    async fn test_login_token_response_shape() {
        let password = "password123";
        let hash = bcrypt::hash(password, bcrypt::DEFAULT_COST).unwrap();
        let user_repo = MockUserRepository::with_user(1, &hash);
        let token_repo = MockTokenRepository::new();

        let input = LoginInput {
            email: "test@example.com".to_string(),
            password: password.to_string(),
        };

        let result = login(&user_repo, &token_repo, "secret", input).await.unwrap();

        // token_type は常に "Bearer"
        assert_eq!(result.token_type, "Bearer");
        // expires_in は 3600 秒固定
        assert_eq!(result.expires_in, 3600);
        // access_token と refresh_token は異なる値
        assert_ne!(result.access_token, result.refresh_token);
        // expires_at は現在時刻より後
        assert!(result.expires_at > Utc::now());
    }
}
