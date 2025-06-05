package api

import (
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/genre"
	"filmogophery/internal/app/features/health"
	"filmogophery/internal/app/features/media"
	"filmogophery/internal/app/features/movie"
)

func RegisterV0Routes(e *echo.Echo, m ...echo.MiddlewareFunc) {
	g := e.Group("v0", m...)

	// --- Handler --- //
	g.GET("/health", health.BuildCheckHealthHandler())

	// --- Movie --- //
	g.GET("/movies", movie.BuildMockedGetMoviesHandler())
	g.GET("/movies/:id", movie.BuildMockedGetMovieDetailHandler())
	g.POST("/movies/:id/impression", movie.BuildMockedPostMovieImpressionHandler())
	g.PUT("/movies/:id/impression", movie.BuildMockedPutMovieImpressionHandler())

	// --- Master --- //
	g.GET("/genres", genre.BuildMockedGetGenresHandler())
	g.GET("/media", media.BuildMockedGetMediaHandler())
}
