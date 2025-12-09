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
	putReviewHandler struct {
		interactor review.UpdateReviewUseCase
	}
	putReviewInput struct {
		ReviewID int32    `param:"id"`
		Rating   *float64 `json:"rating"`
		Comment  *string  `json:"comment"`
	}
)

func NewPutReviewHandler(svc services.IServiceContainer) routers.IRoute {
	return &putReviewHandler{
		interactor: review.NewUpdateReviewInteractor(
			svc.ReviewService(),
		),
	}
}

func (h *putReviewHandler) Register(g *echo.Group) {
	g.PUT("/reviews/:id", h.handle)
}

func (h *putReviewHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed PUT review")

	var req putReviewInput
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := req.validate(); err != nil {
		return err
	}
	logger.Info().Msg("successfully validated params")

	err := h.interactor.Run(
		c.Request().Context(),
		req.ReviewID,
		req.Rating,
		req.Comment,
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (req putReviewInput) validate() error {
	if req.Rating == nil && req.Comment == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "both rating and comment cannot be null")
	}

	if req.Rating != nil {
		rating := *req.Rating
		if rating < 0.1 {
			return echo.NewHTTPError(http.StatusBadRequest, "rating must be at least 0.1")
		}
		if 5.0 < rating {
			return echo.NewHTTPError(http.StatusBadRequest, "rating must be at most 5.0")
		}
	}

	return nil
}
