# filmogophery

[![Dependency](https://github.com/takeYY/filmogophery/actions/workflows/dependabot/update-graph/badge.svg)](https://github.com/takeYY/filmogophery/actions/workflows/dependabot/update-graph)
[![Backend Test](https://github.com/takeYY/filmogophery/actions/workflows/backend-test.yml/badge.svg)](https://github.com/takeYY/filmogophery/actions/workflows/backend-test.yml)
[![GitHub Pages](https://github.com/takeYY/filmogophery/actions/workflows/pages/pages-build-deployment/badge.svg)](https://github.com/takeYY/filmogophery/actions/workflows/pages/pages-build-deployment)

## Frontend

- NextJS

## Backend

- Echo
- Hono
- Axum

## How to use it

```
# Run Echo
make up TARGET_COMPOSE=compose.echo.yml

# Run Hono
make up TARGET_COMPOSE=compose.hono.yml

# Run Axum
make up TARGET_COMPOSE=compose.axum.yml
```

## Comparison

| Aspect          | api-echo (Go)           | api-hono (TypeScript)    | api-axum (Rust)          |
| --------------- | ----------------------- | ------------------------ | ------------------------ |
| DI              | uber/fx                 | manual import            | AppState (Clone)         |
| ORM             | GORM + gen              | Drizzle ORM              | sqlx                     |
| Error handling  | echo.HTTPError          | neverthrow / Result type | thiserror + IntoResponse |
| Validation      | go-playground/validator | zod                      | validator                |
| Logger          | zerolog                 | pino                     | tracing                  |
| Auth middleware | echo.MiddlewareFunc     | createMiddleware         | tower Layer              |
| JWT             | golang-jwt/jwt          | hono/jwt                 | jsonwebtoken             |
