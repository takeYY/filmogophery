package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/features/user"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/routers"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/validators"
	"filmogophery/internal/pkg/logger"
)

type (
	createUserHandler struct {
		interactor user.CreateUserUseCase
	}
	createUserRequest struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,password"`
	}
)

func NewCreateUserHandler(svc services.IServiceContainer) routers.IRoute {
	return &createUserHandler{
		interactor: user.NewCreateUserInteractor(
			svc.UserService(),
		),
	}
}

func (h *createUserHandler) Register(g *echo.Group) {
	g.POST("/users", h.handle)
}

func (h *createUserHandler) handle(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("accessed POST create user")

	var req createUserRequest
	if err := c.Bind(&req); err != nil {
		return responses.ParseBindError(err)
	}
	if errs := validators.ValidateRequest(&req); len(errs) > 0 {
		return responses.ValidationError(errs)
	}
	logger.Info().Msg("successfully validated params")

	result, err := h.interactor.Run(
		c.Request().Context(),
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, result)
}

func (req *createUserRequest) Validate() map[string][]string {
	v := validator.New()
	v.RegisterValidation("password", validators.ValidatePassword)
	return validators.StructToErrors(v.Struct(req))
}
