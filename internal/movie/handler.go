package movie

import (
	"context"
	"net/http"
	"strconv"

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
	e.GET("/movies", h.ReaderHandler.GetMovies)
	e.GET("/movies/:id", h.ReaderHandler.GetMovieById)
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
