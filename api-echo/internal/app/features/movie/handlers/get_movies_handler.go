package handlers

import (
	"fmt"
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
		Genre  string `query:"genre"`
		Limit  int32  `query:"limit"`
		Offset int32  `query:"offset"`
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
	if ng := req.validate(); ng != nil {
		return c.String(http.StatusBadRequest, ng.Error())
	}
	logger.Info().Msg("successfully validated params")

	result, err := h.interactor.Run(
		c.Request().Context(),
		req.Genre,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (req *getMoviesInput) validate() error {
	const maxLimit int32 = 12

	if req.Limit == 0 {
		req.Limit = maxLimit
	} else if req.Limit < 0 || req.Limit > maxLimit {
		return fmt.Errorf("limit must be between 1 and %d", maxLimit)
	}

	if req.Offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}

	return nil
}
