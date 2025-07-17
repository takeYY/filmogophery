package health

import (
	"filmogophery/internal/app/routers"
	"filmogophery/internal/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	healthHandler struct {
	}
)

func NewHealthHandler() routers.IRoute {
	return &healthHandler{}
}

func (h *healthHandler) Register(g *echo.Group) {
	g.GET("/health", h.handle)
}

func (h *healthHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET health check")

	return c.NoContent(http.StatusOK)
}

func BuildCheckHealthHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		logger := logger.GetLogger()
		logger.Info().Msg("accessed GET health check")

		return c.NoContent(http.StatusOK)
	}
}
