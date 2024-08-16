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
	WriterHandler struct {
		commandService *CommandService
	}

	handler struct {
		ReaderHandler ReaderHandler
		WriterHandler WriterHandler
	}
)

func NewHandler(queryService *QueryService, commandService *CommandService) *handler {
	return &handler{
		ReaderHandler: ReaderHandler{
			queryService: queryService,
		},
		WriterHandler: WriterHandler{
			commandService: commandService,
		},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	// Read
	e.GET("/movie/records", h.ReaderHandler.GetMovieWatchRecords)
	// Create
	e.POST("/movie/record", h.WriterHandler.CreateRecord)
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

func (wh *WriterHandler) CreateRecord(c echo.Context) error {
	var dto CreateMovieRecordDto
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: err.Error(),
		})
	}

	result := wh.commandService.CreateRecord(&dto)
	if result != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: result.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response.OK{
		Message: "movie record is created",
	})
}
