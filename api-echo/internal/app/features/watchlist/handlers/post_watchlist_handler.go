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
	postWatchlistHandler struct {
		interactor watchlist.AddWatchlistUseCase
	}
	postWatchlistInput struct {
		MovieID  int32 `json:"movieId" validate:"required"`
		Priority int32 `json:"priority"`
	}
)

func NewPostWatchlistHandler(
	movieRepo repositories.IMovieRepository,
	watchlistRepo repositories.IWatchlistRepository,
) routers.IRoute {
	return &postWatchlistHandler{
		interactor: watchlist.NewAddWatchlistInteractor(
			movieRepo, watchlistRepo,
		),
	}
}

func (h *postWatchlistHandler) RequireAuth() bool {
	return true
}

func (h *postWatchlistHandler) Register(g *echo.Group) {
	g.POST("/watchlist", h.handle)
}

func (h *postWatchlistHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed POST watchlist")

	var req postWatchlistInput
	if err := c.Bind(&req); err != nil {
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
	}
	logger.Info().Msg("successfully validated params")

	err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
		req.MovieID,
		req.Priority,
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (req *postWatchlistInput) Validate() map[string][]string {
	// デフォルト値の設定
	if req.Priority == 0 {
		req.Priority = 1
	}

	v := validator.New()
	return validators.StructToErrors(v.Struct(req))
}
