package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"filmogophery/internal/app/features/platform"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
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

func (h *getPlatformsHandler) RequireAuth() bool {
	return true
}

func (h *getPlatformsHandler) Register(g *echo.Group) {
	g.GET("/platforms", h.handle)
}

func (h *getPlatformsHandler) handle(c echo.Context) error {
	log := zerolog.Ctx(c.Request().Context())
	log.Info().Msg("accessed GET platforms")

	result, err := h.interactor.Run(
		c.Request().Context(),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
