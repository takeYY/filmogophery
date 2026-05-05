# Docs

## Build Redoc

```bash
npx @redocly/cli build-docs docs/openapi.yaml --output docs/index.html
```

## Feature

| OperationID             | Tag          | Echo(Test) | Hono |
| :---------------------- | ------------ | :--------: | :--: |
| health                  | Health       |   ✅(✅)   |  ✅  |
| createUser              | User         |   ✅(🚧)   |  ✅  |
| getCurrentUser          | User         |   ✅( )    |  ✅  |
| getUser                 | User         |    ( )     |      |
| getUserPoints           | User         |   ✅( )    |  ✅  |
| login                   | Auth         |   ✅( )    |  ✅  |
| logout                  | Auth         |   ✅( )    |  ✅  |
| getMovies               | Movie        |   ✅(✅)   |  ✅  |
| getMovieDetail          | Movie        |   ✅(🚧)   |  ✅  |
| searchMovies            | Movie        |   ✅(🚧)   |  ✅  |
| getTrendingMovies       | Trending     |   ✅(🚧)   |  ✅  |
| getMyReviews            | Review       |    ( )     |      |
| createReview            | Review       |   ✅(✅)   |  ✅  |
| updateReview            | Review       |   ✅(🚧)   |  ✅  |
| getMovieWatchHistory    | WatchHistory |   ✅(✅)   |  ✅  |
| createMovieWatchHistory | WatchHistory |   ✅(✅)   |  ✅  |
| getWatchHistory         | WatchHistory |   ✅( )    |  ✅  |
| getWatchlist            | Watchlist    |   ✅(🚧)   |  ✅  |
| addToWatchlist          | Watchlist    |   ✅(🚧)   |  ✅  |
| updateWatchlistItem     | Watchlist    |   🚧(🚧)   |      |
| removeFromWatchlist     | Watchlist    |   ✅(🚧)   |  ✅  |
| getGenres               | Master       |   ✅(✅)   |      |
| getPlatforms            | Master       |   ✅(✅)   |      |
| getSeries               | Master       |   🚧(🚧)   |      |
