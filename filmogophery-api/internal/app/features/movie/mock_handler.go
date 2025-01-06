package movie

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/pkg/logger"
	"filmogophery/internal/pkg/mock"
)

func BuildMockedGetMoviesHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := logger.GetLogger()
		logger.Info().Msg("[Mock] accessed GET movies")

		return c.JSON(http.StatusOK, mock.MockedMovies)
	}
}
