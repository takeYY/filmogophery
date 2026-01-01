package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/movie"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/validators"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	getMovieWatchHistoryHandler struct {
		interactor movie.GetMovieWatchHistoryUseCase
	}
	getMovieWatchHistoryRequest struct {
		MovieID int32 `param:"movieId" validate:"gte=1"`
	}
)

func NewGetMovieWatchHistoryHandler(
	svc services.IServiceContainer,
) routers.IRoute {
	return &getMovieWatchHistoryHandler{
		interactor: movie.NewGetMovieWatchHistoryInteractor(
			svc.MovieService(),
			svc.WatchHistoryService(),
		),
	}
}

func (h *getMovieWatchHistoryHandler) RequireAuth() bool {
	return true
}

func (h *getMovieWatchHistoryHandler) Register(g *echo.Group) {
	g.GET("/movies/:movieId/watch-history", h.handle)
}

func (h *getMovieWatchHistoryHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET movie watch history")

	var req getMovieWatchHistoryRequest
	if err := c.Bind(&req); err != nil {
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
	}
	logger.Info().Msg("successfully validated params")

	result, err := h.interactor.Run(
		c.Request().Context(),
		c.Get("user").(*model.Users),
		req.MovieID,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
func (req *getMovieWatchHistoryRequest) Validate() map[string][]string {
	v := validator.New()
	return validators.StructToErrors(v.Struct(req))
}
