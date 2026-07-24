use std::sync::Arc;

use axum::{Router, middleware};
use sqlx::MySqlPool;
use tower_http::cors::CorsLayer;
use tower_http::trace::TraceLayer;

use crate::app::features::{
    auth,
    health,
    master,
    movie,
    review,
    search,
    trending,
    user,
    watchlist,
};
use crate::config::Config;
use crate::pkg::middleware::auth::require_auth;
use crate::pkg::redis::RedisClient;
use crate::pkg::tmdb::TmdbClient;

/// ハンドラが共有するアプリケーション状態
#[derive(Clone)]
pub struct AppState {
    pub config: Arc<Config>,
    /// DB接続プール。起動時に接続できなかった場合は None になる。
    pub db: Option<MySqlPool>,
    /// TMDB API クライアント
    pub tmdb: Arc<TmdbClient>,
    /// Redis クライアント。起動時に接続できなかった場合は None になる。
    pub redis: Option<Arc<RedisClient>>,
}

/// アプリケーションルーターを構築する
pub fn create_router(
    config: Arc<Config>,
    db: Option<MySqlPool>,
    redis: Option<Arc<RedisClient>>,
) -> Router {
    let tmdb = Arc::new(TmdbClient::new(&config.tmdb.access_token));

    let state = AppState {
        config: Arc::clone(&config),
        db,
        tmdb,
        redis,
    };

    // 認証不要ルート
    let public_routes = Router::new()
        .merge(health::routes())
        .merge(auth::public_routes())
        .merge(user::public_routes());

    // 認証必要ルート
    let protected_routes = Router::new()
        .merge(auth::protected_routes())
        .merge(user::protected_routes())
        .merge(master::routes())
        .merge(movie::routes())
        .merge(review::routes())
        .merge(watchlist::routes())
        .merge(trending::routes())
        .merge(search::routes())
        .layer(middleware::from_fn_with_state(
            Arc::clone(&config),
            require_auth,
        ));

    Router::new()
        .nest(
            "/v1",
            Router::new()
                .merge(public_routes)
                .merge(protected_routes),
        )
        .layer(CorsLayer::permissive())
        .layer(TraceLayer::new_for_http())
        .with_state(state)
}
