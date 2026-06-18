package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"filmogophery/internal/app/routers"
)

type (
	checkHealthHandler struct {
	}
)

func NewCheckHealthHandler() routers.IRoute {
	return &checkHealthHandler{}
}

func (h *checkHealthHandler) RequireAuth() bool {
	return false
}

func (h *checkHealthHandler) Register(g *echo.Group) {
	g.GET("/health", h.handle)
}

func (h *checkHealthHandler) handle(c echo.Context) error {
	log := zerolog.Ctx(c.Request().Context())
	log.Info().Msg("accessed GET health")

	return c.JSON(http.StatusOK, "system all green")
}
