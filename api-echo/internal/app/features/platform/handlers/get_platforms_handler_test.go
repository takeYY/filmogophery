package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"filmogophery/internal/app/features/platform"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/config"
	"filmogophery/tests"
)

func TestGetPlatformsHandler_handle(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	svc := services.NewServiceContainer(testDB, conf, nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/genres", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &getPlatformsHandler{
		interactor: platform.NewGetPlatformsInteractor(svc.PlatformService()),
	}

	// Act
	err := handler.handle(c)

	// expected
	expected := `[
		{"id": 1,  "code": "primeVideo", "name": "Prime Video"},
		{"id": 2,  "code": "netflix",    "name": "Netflix"},
		{"id": 3,  "code": "uNext",      "name": "U-NEXT"},
		{"id": 4,  "code": "disneyPlus", "name": "Disney+"},
		{"id": 5,  "code": "youtube",    "name": "YouTube"},
		{"id": 6,  "code": "appleTv",    "name": "Apple TV+"},
		{"id": 7,  "code": "hulu",       "name": "Hulu"},
		{"id": 8,  "code": "dAnime",     "name": "dアニメ"},
		{"id": 9,  "code": "telasa",     "name": "TELASA"},
		{"id": 10, "code": "cinema",     "name": "映画館"},
		{"id": 99, "code": "unknown",    "name": "不明"}
	]`

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
