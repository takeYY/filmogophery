package handlers

import (
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/routers"
	"filmogophery/internal/pkg/mock"
)

type (
	mockedGenreHandler struct{}
)

func NewMockedGenreHandler() routers.IRoute {
	return &mockedGenreHandler{}
}

func (h *mockedGenreHandler) Register(g *echo.Group) {
	g.GET("/genres", h.getGenres)
}

func (h *mockedGenreHandler) getGenres(c echo.Context) error {
	return c.JSON(200, mock.MockedGenres)
}
