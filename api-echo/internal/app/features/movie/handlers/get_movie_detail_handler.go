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
	getMovieDetailHandler struct {
		interactor movie.GetMovieDetailsUseCase
	}
	getMovieDetailInput struct {
		ID int32 `param:"id" validate:"gte=1"`
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

func (h *getMovieDetailHandler) RequireAuth() bool {
	return true
}

func (h *getMovieDetailHandler) Register(g *echo.Group) {
	g.GET("/movies/:id", h.handle)
}

func (h *getMovieDetailHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET movie detail")

	var req getMovieDetailInput
	if err := c.Bind(&req); err != nil {
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
	}
	logger.Info().Msg("successfully validated params")

	result, err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
		req.ID,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (req *getMovieDetailInput) Validate() map[string][]string {
	v := validator.New()
	return validators.StructToErrors(v.Struct(req))
}
