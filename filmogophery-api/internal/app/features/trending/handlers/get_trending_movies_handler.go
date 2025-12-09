package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/trending"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

type (
	getTrendingMoviesHandler struct {
		interactor trending.GetTrendingMoviesUseCase
	}
)

func NewGetTrendingMoviesHandler(svc services.IServiceContainer) routers.IRoute {
	return &getTrendingMoviesHandler{
		interactor: trending.NewGetTrendingMoviesInteractor(
			svc.DB(),
			svc.MovieService(),
			svc.RedisService(),
			svc.TmdbService(),
		),
	}
}

func (h *getTrendingMoviesHandler) Register(g *echo.Group) {
	g.GET("/trending/movies", h.handle)
}

func (h *getTrendingMoviesHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET trending movies")

	result, err := h.interactor.Run(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
