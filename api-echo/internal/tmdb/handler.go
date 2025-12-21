package tmdb

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/responses"
)

type (
	TmdbHandler struct {
		TmdbService *TmdbService
	}

	handler struct {
		TmdbHandler TmdbHandler
	}
)

func NewHandler(service *TmdbService) *handler {
	return &handler{
		TmdbHandler: TmdbHandler{
			TmdbService: service,
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
		return responses.BadRequestError(map[string][]string{"q": {"must not be empty"}})
	}

	movies, err := th.TmdbService.SearchMovies(&q)
	if err != nil {
		return responses.NotFoundError("movie", map[string][]string{"q": {q}})
	}

	return c.JSON(http.StatusOK, movies)
}
