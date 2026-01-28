package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/user"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/validators"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	getWatchHistoryHandler struct {
		interactor user.GetWatchHistoryUseCase
	}
	getWatchHistoryRequest struct {
		Limit  int32 `query:"limit"`
		Offset int32 `query:"offset"`
	}
)

func NewGetWatchHistoryHandler(svc services.IServiceContainer) routers.IRoute {
	return &getWatchHistoryHandler{
		interactor: user.NewGetWatchHistoryInteractor(
			svc.WatchHistoryService(),
		),
	}
}

func (h *getWatchHistoryHandler) RequireAuth() bool {
	return true
}

func (h *getWatchHistoryHandler) Register(g *echo.Group) {
	g.GET("/users/me/watch-history", h.handle)
}

func (h *getWatchHistoryHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET user watch history")

	var req getWatchHistoryRequest
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
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (req *getWatchHistoryRequest) Validate() map[string][]string {
	// デフォルト値の設定
	if req.Limit == 0 {
		req.Limit = 12
	}

	v := validator.New()
	return validators.StructToErrors(v.Struct(req))
}
