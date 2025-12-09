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
	getMoviesHandler struct {
		interactor movie.GetMoviesUseCase
	}
	getMoviesInput struct {
		Genre string `query:"genre"`
		Limit int32  `query:"limit"`
	}
)

func NewGetMoviesHandler(svc services.IServiceContainer) routers.IRoute {
	return &getMoviesHandler{
		interactor: movie.NewGetMoviesInteractor(svc.MovieService()),
	}
}

func (h *getMoviesHandler) Register(g *echo.Group) {
	g.GET("/movies", h.handle)
}

func (h *getMoviesHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET movies")

	var req getMoviesInput
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	logger.Info().Msg("successfully validated params")

	if req.Limit == 0 { // デフォルトを設定
		req.Limit = 50
	}

	result, err := h.interactor.Run(
		c.Request().Context(),
		req.Genre,
		req.Limit,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
