package media

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/pkg/response"
)

type (
	ReaderHandler struct {
		queryService *QueryService
	}

	handler struct {
		ReaderHandler ReaderHandler
	}
)

func NewHandler() *handler {
	return &handler{
		ReaderHandler: ReaderHandler{
			queryService: NewQueryService(),
		},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	// Read
	e.GET("/media", h.ReaderHandler.GetMedia)
}

func (rh *ReaderHandler) GetMedia(c echo.Context) error {
	media, err := rh.queryService.GetWatchMedia(context.Background())
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "watch media are not found",
		})
	}

	return c.JSON(http.StatusOK, media)
}
