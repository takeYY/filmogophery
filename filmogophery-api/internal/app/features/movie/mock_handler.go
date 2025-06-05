package movie

import (
	"fmt"
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

func BuildMockedGetMovieDetailHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := logger.GetLogger()
		logger.Info().Msg("[Mock] accessed GET movie detail")

		var req GetMovieDetailRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		logger.Info().Msg("successfully validated")

		result, ok := mock.MockedMovieDetailMapper[req.ID]
		if !ok {
			return c.String(http.StatusNotFound, fmt.Sprintf("movie(id=%d) is not found", req.ID))
		}

		return c.JSON(http.StatusOK, result)
	}
}

func BuildMockedPostMovieImpressionHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := logger.GetLogger()
		logger.Info().Msg("[Mock] accessed POST movie impression")

		var req PostMovieImpression
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		logger.Info().Msg("successfully validated")

		return c.NoContent(http.StatusNoContent)
	}
}

func BuildMockedPutMovieImpressionHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := logger.GetLogger()
		logger.Info().Msg("[Mock] accessed PUT movie impression")

		var req PutMovieImpression
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		logger.Info().Msg("successfully validated")

		return c.NoContent(http.StatusNoContent)
	}
}
