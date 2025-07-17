package movie

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

type (
	getMoviesHandler struct {
		interactor GetMoviesUseCase
	}
)

func NewGetMoviesHandler(svc services.IServiceContainer) routers.IRoute {
	return &getMoviesHandler{
		interactor: NewGetMoviesInteractor(svc.MovieService()),
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

// NOTE: test でしか使っていないから移行したら消すこと
func BuildGetMoviesHandler(
	movieService services.IMovieService,
) func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := logger.GetLogger()
		logger.Info().Msg("accessed GET movies")

		interactor := NewGetMoviesInteractor(movieService)
		result, err := interactor.Run(c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, result)
	}
}
