package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/genre"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

type (
	getGenresHandler struct {
		interactor genre.GetGenresUseCase
	}
)

func NewGetGenresHandler(svc services.IServiceContainer) routers.IRoute {
	return &getGenresHandler{
		interactor: genre.NewGetGenresInteractor(
			svc.GenreService(),
		),
	}
}

func (h *getGenresHandler) RequireAuth() bool {
	return true
}

func (h *getGenresHandler) Register(g *echo.Group) {
	g.GET("/genres", h.handle)
}

func (h *getGenresHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET genres")

	result, err := h.interactor.Run(
		c.Request().Context(),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
