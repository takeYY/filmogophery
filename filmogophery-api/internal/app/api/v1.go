package api

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"filmogophery/internal/app/features/health"
	"filmogophery/internal/app/features/movie"
	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/services"
)

func RegisterV1Routes(e *echo.Echo, gormDB *gorm.DB, m ...echo.MiddlewareFunc) {
	g := e.Group("v1", m...)

	// --- Init Repository --- //
	movieRepo := repositories.NewMovieRepository(gormDB)

	// --- Init Service --- //
	movieService := services.NewMovieService(*movieRepo)

	// --- Handler --- //
	g.GET("/health", health.BuildCheckHealthHandler())

	g.GET("/movies", movie.BuildGetMoviesHandler(movieService))
}
