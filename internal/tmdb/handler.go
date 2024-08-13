package tmdb

import (
	"filmogophery/internal/config"
	"filmogophery/pkg/response"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	TmdbHandler struct {
		TmdbService *TmdbService
	}

	handler struct {
		TmdbHandler TmdbHandler
	}
)

func NewHandler(conf *config.Config) *handler {
	return &handler{
		TmdbHandler: TmdbHandler{
			TmdbService: NewTmdbService(conf),
		},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	// Read
	e.GET("/tmdb/search/movies", h.TmdbHandler.SearchTmdbMovies)
}

func (th *TmdbHandler) SearchTmdbMovies(c echo.Context) error {
	q := c.QueryParam("query")
	if q == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: "invalid query params",
		})
	}

	movies, err := th.TmdbService.SearchMovies(&q)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			Message: fmt.Sprintf("movie is not found: %s", q),
		})
	}

	return c.JSON(http.StatusOK, movies)
}
