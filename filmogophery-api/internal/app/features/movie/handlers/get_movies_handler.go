package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/movie"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

type (
	getMoviesHandler struct {
		interactor movie.GetMoviesUseCase
	}
)

func NewGetMoviesHandler(svc services.IServiceContainer) routers.IRoute {
	return &getMoviesHandler{
		interactor: movie.NewGetMoviesInteractor(svc.MovieService()),
	}
}

func (h *getMoviesHandler) Register(g *echo.Group) {
	g.GET("/movies", h.handle)
}

func (h *getMoviesHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET movies")

	result, err := h.interactor.Run(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
