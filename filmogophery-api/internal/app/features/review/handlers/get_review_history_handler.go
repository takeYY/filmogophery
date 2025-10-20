package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/review"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/logger"
)

type (
	getReviewHistoryHandler struct {
		interactor review.GetReviewHistoryUseCase
	}
	getReviewHistoryInput struct {
		ID int32 `param:"id"`
	}
)

func NewGetReviewHistoryHandler(svc services.IServiceContainer) routers.IRoute {
	return &getReviewHistoryHandler{
		interactor: review.NewGetReviewHistoryInteractor(svc.ReviewService()),
	}
}

func (h *getReviewHistoryHandler) Register(g *echo.Group) {
	g.GET("/reviews/:id/history", h.handle)
}

func (h *getReviewHistoryHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed GET review history")

	var req getReviewHistoryInput
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	logger.Info().Msg("successfully validated params")

	result, err := h.interactor.Run(
		c.Request().Context(),
		req.ID,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
