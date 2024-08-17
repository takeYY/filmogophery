package impression

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/pkg/response"
)

type (
	ReaderHandler struct {
		queryService *QueryService
	}

	handler struct {
		ReaderHandler ReaderHandler
	}
)

func NewHandler(queryService *QueryService) *handler {
	return &handler{
		ReaderHandler: ReaderHandler{
			queryService: queryService,
		},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	// Read
	e.GET("/movie/impressions", h.ReaderHandler.GetImpressions)
}

func (rh *ReaderHandler) GetImpressions(c echo.Context) error {
	impressions, err := rh.queryService.GetImpressions(context.Background())
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "impressions are not found",
		})
	}

	return c.JSON(http.StatusOK, impressions)
}
