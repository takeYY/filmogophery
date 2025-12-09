package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/search"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

type (
	searchMoviesHandler struct {
		interactor search.SearchMoviesUseCase
	}
	searchMoviesInput struct {
		Title  string `query:"title"`
		Limit  *int32 `query:"limit"`
		Offset *int32 `query:"offset"`
	}
)

func NewSearchMoviesHandler(svc services.IServiceContainer) routers.IRoute {
	return &searchMoviesHandler{
		interactor: search.NewSearchMoviesInteractor(
			svc.DB(),
			svc.MovieService(),
			svc.RedisService(),
			svc.TmdbService(),
		),
	}
}

func (h *searchMoviesHandler) Register(g *echo.Group) {
	g.GET("/search/movies", h.handle)
}

func (h *searchMoviesHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET search movies")

	var req searchMoviesInput
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if ng := req.validate(); ng != nil {
		return ng
	}
	logger.Info().Msg("successfully validated params")

	result, err := h.interactor.Run(c.Request().Context(), req.Title, *req.Limit, *req.Offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (req *searchMoviesInput) validate() error {
	if req.Limit == nil {
		req.Limit = &[]int32{20}[0]
	}
	if req.Offset == nil {
		req.Offset = &[]int32{0}[0]
	}

	if req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "title cannot be null")
	}

	if *req.Limit < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "limit must be at least 0")
	}
	if 20 < *req.Limit {
		return echo.NewHTTPError(http.StatusBadRequest, "limit must be at most 20")
	}

	if *req.Offset < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "offset must be at least 0")
	}

	return nil
}
