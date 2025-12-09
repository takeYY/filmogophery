# Docs

## Build Redoc

```bash
npx @redocly/cli build-docs docs/openapi.yaml --output docs/index.html
```

## Feature

| OperationID         | Tag       | Echo | Test |
| :------------------ | --------- | :--: | :--: |
| health              | Health    |  ✅  |  ✅  |
| createUser          | User      |  🚧  |  🚧  |
| getMovies           | Movie     |  ✅  |  ✅  |
| getMovieDetail      | Movie     |  ✅  |  🚧  |
| searchMovies        | Movie     |  ✅  |  🚧  |
| getTrendingMovies   | Trending  |  ✅  |  🚧  |
| createReview        | Review    |  ✅  |  ✅  |
| updateReview        | Review    |  ✅  |  🚧  |
| getWatchHistory     | Review    |  ✅  |  ✅  |
| addWatchHistory     | Review    |  ✅  |  ✅  |
| getWatchlist        | Watchlist |  🚧  |  🚧  |
| addToWatchlist      | Watchlist |  🚧  |  🚧  |
| updateWatchlistItem | Watchlist |  🚧  |  🚧  |
| removeFromWatchlist | Watchlist |  🚧  |  🚧  |
| getGenres           | Master    |  ✅  |  ✅  |
| getPlatforms        | Master    |  ✅  |  ✅  |
| getSeries           | Master    |  🚧  |  🚧  |
