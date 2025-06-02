package genre

import (
	"filmogophery/internal/pkg/mock"

	"github.com/labstack/echo/v4"
)

func BuildMockedGetGenresHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(200, mock.MockedGenres)
	}
}
