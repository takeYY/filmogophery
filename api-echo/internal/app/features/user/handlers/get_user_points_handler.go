package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"filmogophery/internal/app/features/user"
	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/validators"
	"filmogophery/internal/pkg/gen/model"
)

type (
	getUserPointsHandler struct {
		interactor user.GetUserPointsUseCase
	}
	getUserPointsQuery struct {
		Limit  int32 `query:"limit" validate:"gte=1,lte=50"`
		Offset int32 `query:"offset" validate:"gte=0"`
	}
)

func NewGetUserPointsHandler(svc services.IServiceContainer) routers.IRoute {
	return &getUserPointsHandler{
		interactor: user.NewGetUserPointsInteractor(
			svc.PointService(),
			repositories.NewPointRepository(svc.DB()),
		),
	}
}

func (h *getUserPointsHandler) RequireAuth() bool {
	return true
}

func (h *getUserPointsHandler) Register(g *echo.Group) {
	g.GET("/users/me/points", h.handle)
}

func (h *getUserPointsHandler) handle(c echo.Context) error {
	log := zerolog.Ctx(c.Request().Context())
	log.Info().Msg("accessed GET user points")

	req := getUserPointsQuery{Limit: 20, Offset: 0}
	if err := c.Bind(&req); err != nil {
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
	}

	result, err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (req *getUserPointsQuery) Validate() map[string][]string {
	v := validator.New()
	return validators.StructToErrors(v.Struct(req))
}
