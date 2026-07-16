use axum::Router;

use crate::app::router::AppState;

// TODO: POST /v1/movies/:id/reviews              (requires auth)
// TODO: PUT  /v1/movies/:id/reviews              (requires auth)
// TODO: GET  /v1/movies/:id/reviews              (requires auth)
// TODO: POST /v1/movies/:id/reviews/watch-history (requires auth)

pub fn routes() -> Router<AppState> {
    Router::new()
}
