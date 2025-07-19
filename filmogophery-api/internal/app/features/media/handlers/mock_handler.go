package handlers

import (
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/routers"
	"filmogophery/internal/pkg/mock"
)

type (
	mockedMediaHandler struct{}
)

func NewMockedMediaHandler() routers.IRoute {
	return &mockedMediaHandler{}
}

func (h *mockedMediaHandler) Register(g *echo.Group) {
	g.GET("/media", h.getMedia)
}

func (h *mockedMediaHandler) getMedia(c echo.Context) error {
	return c.JSON(200, mock.MockedMedia)
}
