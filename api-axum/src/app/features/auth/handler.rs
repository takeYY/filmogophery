use axum::Router;

use crate::app::router::AppState;

// TODO: POST /v1/auth/login
// TODO: POST /v1/auth/logout  (requires auth)

pub fn routes() -> Router<AppState> {
    Router::new()
}
