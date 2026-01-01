# Docs

## Build Redoc

```bash
npx @redocly/cli build-docs docs/openapi.yaml --output docs/index.html
```

## Feature

| OperationID             | Tag          | Echo | Test |
| :---------------------- | ------------ | :--: | :--: |
| health                  | Health       |  ✅  |  ✅  |
| createUser              | User         |  ✅  |  🚧  |
| getMyUser               | User         |      |      |
| getUser                 | User         |      |      |
| login                   | Auth         |  ✅  |      |
| logout                  | Auth         |      |      |
| getMovies               | Movie        |  ✅  |  ✅  |
| getMovieDetail          | Movie        |  ✅  |  🚧  |
| searchMovies            | Movie        |  ✅  |  🚧  |
| getTrendingMovies       | Trending     |  ✅  |  🚧  |
| getMyReviews            | Review       |      |      |
| createReview            | Review       |  ✅  |  ✅  |
| updateReview            | Review       |  ✅  |  🚧  |
| getMovieWatchHistory    | WatchHistory |  ✅  |  ✅  |
| createMovieWatchHistory | WatchHistory |  ✅  |  ✅  |
| getWatchlist            | Watchlist    |  ✅  |  🚧  |
| addToWatchlist          | Watchlist    |  🚧  |  🚧  |
| updateWatchlistItem     | Watchlist    |  🚧  |  🚧  |
| removeFromWatchlist     | Watchlist    |  🚧  |  🚧  |
| getGenres               | Master       |  ✅  |  ✅  |
| getPlatforms            | Master       |  ✅  |  ✅  |
| getSeries               | Master       |  🚧  |  🚧  |
