package health

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/config"
	"filmogophery/internal/pkg/logger"
	"filmogophery/internal/pkg/response"
)

type (
	ReaderHandler struct{}
	handler       struct {
		ReaderHandler ReaderHandler
	}
)

func NewHandler(conf *config.Config) *handler {
	return &handler{
		ReaderHandler: ReaderHandler{},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	// Read
	e.GET("/health", h.ReaderHandler.Health)
}

func (h *ReaderHandler) Health(c echo.Context) error {
	logger := logger.GetLogger()

	logger.Info().Msg("accessed health check!!")

	return c.JSON(http.StatusOK, response.OK{
		Message: "system all green",
	})
}
