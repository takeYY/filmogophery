package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/review"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/validators"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	putReviewHandler struct {
		interactor review.UpdateReviewUseCase
	}
	putReviewInput struct {
		ReviewID int32    `param:"id" validate:"gte=1"`
		Rating   *float64 `json:"rating"`  // required_without=Comment が上手く機能しないので Validate() 内で対応
		Comment  *string  `json:"comment"` // required_without=Rating が上手く機能しないので Validate() 内で対応
	}
)

func NewPutReviewHandler(svc services.IServiceContainer) routers.IRoute {
	return &putReviewHandler{
		interactor: review.NewUpdateReviewInteractor(
			svc.ReviewService(),
		),
	}
}

func (h *putReviewHandler) RequireAuth() bool {
	return true
}

func (h *putReviewHandler) Register(g *echo.Group) {
	g.PUT("/reviews/:id", h.handle)
}

func (h *putReviewHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed PUT review")

	var req putReviewInput
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
		req.ReviewID,
		req.Rating,
		req.Comment,
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (req *putReviewInput) Validate() map[string][]string {
	errors := make(map[string][]string)
	if req.Rating != nil {
		rating := *req.Rating
		if rating < 0.1 {
			errors["Rating"] = append(errors["Rating"], "Rating validation failed on gte")
		} else if 5.0 < rating {
			errors["Rating"] = append(errors["Rating"], "Rating validation failed on lte")
		}
	} else if req.Comment == nil {
		errors["Rating"] = append(errors["Rating"], "both rating and comment cannot be null")
		errors["Comment"] = append(errors["Comment"], "both rating and comment cannot be null")
	}

	return errors
}
