package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"filmogophery/internal/app/features/user"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/pkg/gen/model"
)

type (
	getCurrentUserHandler struct {
		interactor user.GetCurrentUserUseCase
	}
)

func NewGetCurrentUserHandler() routers.IRoute {
	return &getCurrentUserHandler{
		interactor: user.NewGetCurrentUserInteractor(),
	}
}

func (h *getCurrentUserHandler) RequireAuth() bool {
	return true
}

func (h *getCurrentUserHandler) Register(g *echo.Group) {
	g.GET("/users/me", h.handle)
}

func (h *getCurrentUserHandler) handle(c echo.Context) error {
	log := zerolog.Ctx(c.Request().Context())
	log.Info().Msg("accessed GET current user")

	result, err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
