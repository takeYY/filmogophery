package api

import (
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/health"
	"filmogophery/internal/app/features/movie"
)

func RegisterV0Routes(e *echo.Echo, m ...echo.MiddlewareFunc) {
	g := e.Group("v0", m...)

	// --- Handler --- //
	g.GET("/health", health.BuildCheckHealthHandler())

	g.GET("/movies", movie.BuildMockedGetMoviesHandler())
}
