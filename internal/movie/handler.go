package movie

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"filmogophery/pkg/response"
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

func NewHandler() *handler {
	return &handler{
		ReaderHandler: ReaderHandler{
			queryService: NewQueryService(),
		},
		WriterHandler: WriterHandler{
			commandService: NewCommandService(),
		},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	// Read
	e.GET("/movies", h.ReaderHandler.GetMovies)
	e.GET("/movies/:id", h.ReaderHandler.GetMovieById)
	// Create
	e.POST("movie", h.WriterHandler.Create)
}

func (rh *ReaderHandler) GetMovieById(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "invalid movie id",
		})
	}

	movie, err := rh.queryService.GetMovieDetails(context.Background(), &id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "movie is not found",
		})
	}

	return c.JSON(http.StatusOK, movie)
}

func (rh *ReaderHandler) GetMovies(c echo.Context) error {
	movies, err := rh.queryService.GetMovies(context.Background())
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: "movies are not found",
		})
	}

	return c.JSON(http.StatusOK, movies)
}

func (wh *WriterHandler) Create(c echo.Context) error {
	var dto CreateMovieDto
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: err.Error(),
		})
	}

	movie, err := wh.commandService.CreateMovie(&dto)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: fmt.Sprintf("movie can not create: %v", err),
		})
	}

	return c.JSON(http.StatusOK, movie)
}
