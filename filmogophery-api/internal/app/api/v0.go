package api

import (
	"filmogophery/internal/app/features/health"

	"github.com/labstack/echo/v4"
)

func RegisterV0Routes(e *echo.Echo, m ...echo.MiddlewareFunc) {
	g := e.Group("v0", m...)

	// --- Handler --- //
	g.GET("/health", health.BuildCheckHealthHandler())
}
