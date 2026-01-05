package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/auth"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	logoutHandler struct {
		interactor auth.LogoutUseCase
	}
)

func NewLogoutHandler(svc services.IServiceContainer) routers.IRoute {
	return &logoutHandler{
		interactor: auth.NewLogoutInteractor(
			svc.UserService(),
		),
	}
}

func (h *logoutHandler) RequireAuth() bool {
	return true
}

func (h *logoutHandler) Register(g *echo.Group) {
	g.POST("/auth/logout", h.handle)
}

func (h *logoutHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed POST logout")

	err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
