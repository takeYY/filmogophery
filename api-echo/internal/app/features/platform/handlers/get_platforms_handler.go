package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/platform"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

type (
	getPlatformsHandler struct {
		interactor platform.GetPlatformsUseCase
	}
)

func NewGetPlatformsHandler(svc services.IServiceContainer) routers.IRoute {
	return &getPlatformsHandler{
		interactor: platform.NewGetPlatformsInteractor(
			svc.PlatformService(),
		),
	}
}

func (h *getPlatformsHandler) Register(g *echo.Group) {
	g.GET("/platforms", h.handle)
}

func (h *getPlatformsHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET platforms")

	result, err := h.interactor.Run(
		c.Request().Context(),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
