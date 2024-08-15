package impression

import (
	"context"
	"filmogophery/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
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
