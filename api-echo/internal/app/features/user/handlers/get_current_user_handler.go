package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/user"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
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
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET current user")

	result, err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
