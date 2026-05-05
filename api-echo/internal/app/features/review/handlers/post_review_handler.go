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
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	postReviewHandler struct {
		interactor review.CreateReviewUseCase
	}
	postReviewInput struct {
		MovieID      int32              `param:"id" validate:"gte=1"`
		Rating       *float64           `json:"rating"`
		Comment      *string            `json:"comment"`
		WatchHistory *watchHistoryInput `json:"watchHistory"`
	}
	watchHistoryInput struct {
		PlatformID  int32          `json:"platformId" validate:"required,gte=1"`
		WatchedDate *constant.Date `json:"watchedDate"`
	}
)

func NewPostReviewHandler(svc services.IServiceContainer) routers.IRoute {
	return &postReviewHandler{
		interactor: review.NewCreateReviewInteractor(
			svc.DB(),
			svc.MovieService(),
			svc.ReviewService(),
			svc.PlatformService(),
			svc.PointService(),
		),
	}
}

func (h *postReviewHandler) RequireAuth() bool {
	return true
}

func (h *postReviewHandler) Register(g *echo.Group) {
	g.POST("/movies/:id/reviews", h.handle)
}

func (h *postReviewHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed POST review")

	var req postReviewInput
	if err := c.Bind(&req); err != nil {
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
	}
	logger.Info().Msg("successfully validated params")

	var whInput *review.WatchHistoryInput
	if req.WatchHistory != nil {
		whInput = &review.WatchHistoryInput{
			PlatformID:  req.WatchHistory.PlatformID,
			WatchedDate: req.WatchHistory.WatchedDate,
		}
	}

	err := h.interactor.Run(
		c.Request().Context(),
		c.Get("operator").(*model.Users),
		req.MovieID,
		req.Rating,
		req.Comment,
		whInput,
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (req *postReviewInput) Validate() map[string][]string {
	errors := make(map[string][]string)

	// rating/comment のどちらかは必須
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

	// 視聴履歴のバリデーション
	if req.WatchHistory != nil {
		v := validator.New()
		for k, v := range validators.StructToErrors(v.Struct(req.WatchHistory)) {
			errors[k] = v
		}
		if req.WatchHistory.WatchedDate != nil {
			parsedTime, err := time.ParseInLocation(constant.DateFormat, string(*req.WatchHistory.WatchedDate), time.Local)
			if err != nil {
				errors["WatchedDate"] = append(errors["WatchedDate"], fmt.Sprintf("failed to parse date(%s)", *req.WatchHistory.WatchedDate))
			} else if parsedTime.After(time.Now()) {
				errors["WatchedDate"] = append(errors["WatchedDate"], "date cannot be in the future")
			}
		}
	}

	return errors
}
