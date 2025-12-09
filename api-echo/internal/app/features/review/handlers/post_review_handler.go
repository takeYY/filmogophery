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
	postReviewHandler struct {
		interactor review.CreateReviewUseCase
	}
	postReviewInput struct {
		MovieID int32    `param:"id"`
		Rating  *float64 `json:"rating"`
		Comment *string  `json:"comment"`
	}
)

func NewPostReviewHandler(svc services.IServiceContainer) routers.IRoute {
	return &postReviewHandler{
		interactor: review.NewCreateReviewInteractor(
			svc.MovieService(),
			svc.ReviewService(),
		),
	}
}

func (h *postReviewHandler) Register(g *echo.Group) {
	g.POST("/movies/:id/reviews", h.handle)
}

func (h *postReviewHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed POST review")

	var req postReviewInput
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := req.validate(); err != nil {
		return err
	}
	logger.Info().Msg("successfully validated params")

	err := h.interactor.Run(
		c.Request().Context(),
		req.MovieID,
		req.Rating,
		req.Comment,
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (req postReviewInput) validate() error {
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
