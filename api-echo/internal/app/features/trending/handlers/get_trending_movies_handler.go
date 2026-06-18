package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"filmogophery/internal/app/features/trending"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/gen/model"
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
			svc.ReviewService(),
			svc.RedisService(),
			svc.TmdbService(),
		),
	}
}

func (h *getTrendingMoviesHandler) RequireAuth() bool {
	return true
}

func (h *getTrendingMoviesHandler) Register(g *echo.Group) {
	g.GET("/trending/movies", h.handle)
}

func (h *getTrendingMoviesHandler) handle(c echo.Context) error {
	log := zerolog.Ctx(c.Request().Context())
	log.Info().Msg("accessed GET trending movies")

	result, err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
