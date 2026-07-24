use serde::Deserialize;

use crate::app::responses::AppError;

const TMDB_BASE_URL: &str = "https://api.themoviedb.org/3";

// ─── TMDB レスポンス型 ────────────────────────────────────────

#[derive(Debug, Deserialize)]
pub struct TmdbSearchResult {
    pub results: Vec<TmdbMovieResult>,
}

#[derive(Debug, Deserialize)]
pub struct TmdbMovieResult {
    pub id: i32,
    pub title: String,
    #[serde(default)]
    pub overview: String,
    #[serde(default)]
    pub release_date: String,
    pub poster_path: Option<String>,
    #[serde(default)]
    pub genre_ids: Vec<i32>,
}

#[derive(Debug, Deserialize)]
pub struct TmdbMovieDetail {
    pub runtime: Option<i32>,
    pub vote_average: f64,
    pub vote_count: i32,
}

impl TmdbMovieDetail {
    /// TMDB の 10 点満点を 5 点満点に変換（小数第1位に丸め）
    pub fn vote_average_5(&self) -> f64 {
        (self.vote_average / 2.0 * 10.0).round() / 10.0
    }
}

// ─── TMDB クライアント ────────────────────────────────────────

pub struct TmdbClient {
    http: reqwest::Client,
    access_token: String,
}

impl TmdbClient {
    pub fn new(access_token: &str) -> Self {
        Self {
            http: reqwest::Client::new(),
            access_token: access_token.to_string(),
        }
    }

    /// TMDB でタイトルによる映画検索（日本語対応）
    pub async fn search_movies(&self, title: &str, page: u32) -> Result<TmdbSearchResult, AppError> {
        let url = format!(
            "{}/search/movie?query={}&language=ja-JP&page={}",
            TMDB_BASE_URL,
            urlencoding::encode(title),
            page
        );

        self.get::<TmdbSearchResult>(&url).await
    }

    /// TMDB から映画詳細を取得
    pub async fn get_movie_detail(&self, tmdb_id: i32) -> Result<TmdbMovieDetail, AppError> {
        let url = format!("{}/movie/{}?language=ja-JP", TMDB_BASE_URL, tmdb_id);
        self.get::<TmdbMovieDetail>(&url).await
    }

    async fn get<T: serde::de::DeserializeOwned>(&self, url: &str) -> Result<T, AppError> {
        self.http
            .get(url)
            .header("Authorization", format!("Bearer {}", self.access_token))
            .header("Accept", "application/json")
            .send()
            .await
            .map_err(|e| {
                tracing::error!("TMDB request failed: {e}");
                AppError::InternalServerError
            })?
            .json::<T>()
            .await
            .map_err(|e| {
                tracing::error!("TMDB response parse failed: {e}");
                AppError::InternalServerError
            })
    }
}
