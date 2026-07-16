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

/// ハンドラが共有するアプリケーション状態
#[derive(Clone)]
pub struct AppState {
    pub config: Arc<Config>,
    /// DB接続プール。起動時に接続できなかった場合は None になる。
    /// 実際のDBアクセスが必要なハンドラでは None チェックを行うこと。
    pub db: Option<MySqlPool>,
}

/// アプリケーションルーターを構築する
pub fn create_router(config: Arc<Config>, db: Option<MySqlPool>) -> Router {
    let state = AppState {
        config: Arc::clone(&config),
        db,
    };

    // 認証不要ルート
    let public_routes = Router::new()
        .merge(health::routes())
        .merge(auth::routes())
        .merge(user::public_routes())
        .merge(master::routes());

    // 認証必要ルート
    let protected_routes = Router::new()
        .merge(user::protected_routes())
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
