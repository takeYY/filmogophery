package routers

import "github.com/labstack/echo/v4"

type (
	IRoute interface {
		Register(g *echo.Group)
		// RequireUser() bool
		// handle(c echo.Context) error
	}
)
