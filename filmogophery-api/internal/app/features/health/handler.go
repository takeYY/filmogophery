package health

import (
	"filmogophery/internal/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

func BuildCheckHealthHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := logger.GetLogger()
		logger.Info().Msg("accessed GET health check")

		return c.NoContent(http.StatusOK)
	}
}
