package record

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
	e.GET("/movie/records", h.ReaderHandler.GetMovieWatchRecords)
}

func (rh *ReaderHandler) GetMovieWatchRecords(c echo.Context) error {
	watchRecords, err := rh.queryService.GetWatchRecords(context.Background())
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "movie watch records are not found",
		})
	}

	return c.JSON(http.StatusOK, watchRecords)
}
