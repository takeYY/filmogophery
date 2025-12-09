package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/movie"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

type (
	getMovieDetailHandler struct {
		interactor movie.GetMovieDetailsUseCase
	}
	getMovieDetailInput struct {
		ID int32 `param:"id"`
	}
)

func NewGetMovieDetailHandler(svc services.IServiceContainer) routers.IRoute {
	return &getMovieDetailHandler{
		interactor: movie.NewGetMovieDetailInteractor(
			svc.MovieService(),
			svc.ReviewService(),
			svc.TmdbService(),
		),
	}
}

func (h *getMovieDetailHandler) Register(g *echo.Group) {
	g.GET("/movies/:id", h.handle)
}

func (h *getMovieDetailHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET movie detail")

	var req getMovieDetailInput
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	logger.Info().Msg("successfully validated params")

	result, err := h.interactor.Run(
		c.Request().Context(),
		req.ID,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
