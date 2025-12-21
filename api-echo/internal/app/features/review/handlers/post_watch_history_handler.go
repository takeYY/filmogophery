package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/review"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/validators"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/logger"
)

type (
	postReviewHistoryHandler struct {
		interactor review.CreateReviewHistoryUseCase
	}
	postReviewHistoryInput struct {
		ReviewID    int32          `param:"id" validate:"gte=1"`
		PlatformID  int32          `json:"platformId" validate:"required,gte=1"`
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
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
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

func (req *postReviewHistoryInput) Validate() map[string][]string {
	errors := make(map[string][]string)

	if req.WatchedDate != nil {
		parsedTime, err := time.ParseInLocation(constant.DateFormat, string(*req.WatchedDate), time.Local)
		if err != nil {
			errors["WatchedDate"] = append(errors["WatchedDate"], fmt.Sprintf("failed to parse date(%s)", *req.WatchedDate))
			return errors
		}

		now := time.Now()
		if parsedTime.After(now) {
			errors["WatchedDate"] = append(errors["WatchedDate"], "date cannot be in the future")
			return errors
		}
	}

	v := validator.New()
	return validators.StructToErrors(v.Struct(req))
}
