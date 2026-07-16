use axum::Router;

use crate::app::router::AppState;

// TODO: POST /v1/users           (public)
// TODO: GET  /v1/users/me        (requires auth)
// TODO: GET  /v1/users/me/watch-history (requires auth)
// TODO: GET  /v1/users/me/points (requires auth)

pub fn public_routes() -> Router<AppState> {
    Router::new()
}

pub fn protected_routes() -> Router<AppState> {
    Router::new()
}
