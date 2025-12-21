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
	"filmogophery/internal/pkg/logger"
)

type (
	getMoviesHandler struct {
		interactor movie.GetMoviesUseCase
	}
	getMoviesInput struct {
		Genre  string `query:"genre"`
		Limit  int32  `query:"limit" validate:"gte=1,lte=12"`
		Offset int32  `query:"offset" validate:"gte=0"`
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
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
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

func (req *getMoviesInput) Validate() map[string][]string {
	// デフォルト値の設定
	if req.Limit == 0 {
		req.Limit = 12
	}

	v := validator.New()
	return validators.StructToErrors(v.Struct(req))
}
