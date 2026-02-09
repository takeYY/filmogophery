package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/watchlist"
	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	deleteWatchlistHandler struct {
		interactor watchlist.DeleteFromWatchlistUseCase
	}
	deleteWatchlistInput struct {
		watchlistID int32 `param:"watchlistId"`
	}
)

func NewDeleteWatchlistHandler(
	watchlistRepo repositories.IWatchlistRepository,
) routers.IRoute {
	return &deleteWatchlistHandler{
		interactor: watchlist.NewDeleteFromWatchlistInteractor(
			watchlistRepo,
		),
	}
}

func (h *deleteWatchlistHandler) RequireAuth() bool {
	return true
}

func (h *deleteWatchlistHandler) Register(g *echo.Group) {
	g.DELETE("/watchlist/:watchlistId", h.handle)
}

func (h *deleteWatchlistHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed DELETE watchlist")

	var req deleteWatchlistInput
	if err := c.Bind(&req); err != nil {
		return responses.ParseBindError(err)
	}
	logger.Info().Msg("successfully validated params")

	err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
		req.watchlistID,
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
