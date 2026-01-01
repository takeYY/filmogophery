package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/auth"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/validators"
	"filmogophery/internal/pkg/logger"
)

type (
	loginHandler struct {
		interactor auth.LoginUseCase
	}
	loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,password"`
	}
)

func NewLoginHandler(svc services.IServiceContainer) routers.IRoute {
	return &loginHandler{
		interactor: auth.NewLoginInteractor(
			svc.UserService(),
		),
	}
}

func (h *loginHandler) Register(g *echo.Group) {
	g.POST("/auth/login", h.handle)
}

func (h *loginHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed POST login")

	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
	}
	logger.Info().Msg("successfully validated params")

	result, err := h.interactor.Run(
		c.Request().Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (req *loginRequest) Validate() map[string][]string {
	v := validator.New()
	v.RegisterValidation("password", validators.ValidatePassword)
	return validators.StructToErrors(v.Struct(req))
}
