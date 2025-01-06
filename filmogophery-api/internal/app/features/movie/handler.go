package movie

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

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
