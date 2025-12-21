package responses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func ParseBindError(err error) error {
	errors := make(map[string][]string)

	if be, ok := err.(*echo.BindingError); ok {
		errors["body"] = []string{be.Error()}
	} else {
		errors["body"] = []string{err.Error()}
	}

	return ValidationError(errors)
}

func ValidationError(errors map[string][]string) error {
	return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{
		Message: "validation failed",
		Errors:  errors,
	})
}

func BadRequestError(errors map[string][]string) error {
	return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{
		Message: "bad request",
		Errors:  errors,
	})
}

func NotFoundError(resource string, errors map[string][]string) error {
	return echo.NewHTTPError(http.StatusNotFound, ErrorResponse{
		Message: resource + " not found",
		Errors:  errors,
	})
}

func InternalServerError() error {
	return echo.NewHTTPError(http.StatusInternalServerError, ErrorResponse{
		Message: "system error",
	})
}
