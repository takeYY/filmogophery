package api

import (
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/genre"
	"filmogophery/internal/app/features/health"
	"filmogophery/internal/app/features/movie"
)

func RegisterV0Routes(e *echo.Echo, m ...echo.MiddlewareFunc) {
	g := e.Group("v0", m...)

	// --- Handler --- //
	g.GET("/health", health.BuildCheckHealthHandler())

	// --- Movie --- //
	g.GET("/movies", movie.BuildMockedGetMoviesHandler())

	// --- Master --- //
	g.GET("/genres", genre.BuildMockedGetGenresHandler())
}
