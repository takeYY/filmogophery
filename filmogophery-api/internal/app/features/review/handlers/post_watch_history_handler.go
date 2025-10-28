package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/review"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/logger"
)

type (
	postReviewHistoryHandler struct {
		interactor review.CreateReviewHistoryUseCase
	}
	postReviewHistoryInput struct {
		ReviewID    int32          `param:"id"`
		PlatformID  int32          `json:"platformId"`
		WatchedDate *constant.Date `json:"watchedDate"`
	}
)

func NewPostReviewHistoryHandler(svc services.IServiceContainer) routers.IRoute {
	return &postReviewHistoryHandler{
		interactor: review.NewCreateReviewHistoryInteractor(
			svc.ReviewService(),
			svc.PlatformService(),
		),
	}
}

func (h *postReviewHistoryHandler) Register(g *echo.Group) {
	g.POST("/reviews/:id/history", h.handle)
}

func (h *postReviewHistoryHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed POST review history")

	var req postReviewHistoryInput
	if err := c.Bind(&req); err != nil {
		logger.Error().Msgf("failed to bind: %s", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	if ng := req.validate(); ng != nil {
		return ng
	}
	logger.Info().Msg("successfully validated params")

	err := h.interactor.Run(
		c.Request().Context(),
		req.ReviewID,
		req.PlatformID,
		req.WatchedDate,
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (req *postReviewHistoryInput) validate() error {
	logger := logger.GetLogger()

	if req.WatchedDate != nil {
		parsedTime, err := time.ParseInLocation(constant.DateFormat, string(*req.WatchedDate), time.Local)
		if err != nil {
			em := fmt.Sprintf("failed to parse watchedDate(%s)", *req.WatchedDate)
			logger.Error().Msg(em + ":" + err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, em)
		}

		now := time.Now()
		if parsedTime.After(now) {
			em := "watchedDate must not be in the future"
			logger.Warn().Msg(em)
			return echo.NewHTTPError(http.StatusBadRequest, em)
		}
	}

	return nil
}
