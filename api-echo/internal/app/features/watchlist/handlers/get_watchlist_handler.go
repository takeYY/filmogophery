package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/watchlist"
	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/validators"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	getWatchlistHandler struct {
		interactor watchlist.GetWatchlistUseCase
	}

	getWatchlistInput struct {
		Limit  int32 `query:"limit" validate:"gte=1,lte=12"`
		Offset int32 `query:"offset" validate:"gte=0"`
	}
)

func NewGetWatchlistHandler(
	watchlistRepo repositories.IWatchlistRepository,
) routers.IRoute {
	return &getWatchlistHandler{
		interactor: watchlist.NewGetWatchlistInteractor(
			watchlistRepo,
		),
	}
}

func (h *getWatchlistHandler) RequireAuth() bool {
	return true
}

func (h *getWatchlistHandler) Register(g *echo.Group) {
	g.GET("/watchlist", h.handle)
}

func (h *getWatchlistHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET watchlist")

	var req getWatchlistInput
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

func (req *getWatchlistInput) Validate() map[string][]string {
	// デフォルト値の設定
	if req.Limit == 0 {
		req.Limit = 12
	}

	v := validator.New()
	return validators.StructToErrors(v.Struct(req))
}
