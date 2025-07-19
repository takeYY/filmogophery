package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/routers"
	"filmogophery/internal/pkg/logger"
	"filmogophery/internal/pkg/mock"
)

type (
	mockedMovieHandler struct{}

	GetMovieDetailRequest struct {
		ID int32 `param:"id"`
	}

	PostMovieImpression struct {
		ID        int32   `param:"id"`
		WatchDate string  `json:"watchDate"`
		MediaCode string  `json:"mediaCode"`
		Rating    float32 `json:"rating"`
		Note      string  `json:"note"`
	}

	PutMovieImpression struct {
		ID     int32   `param:"id"`
		Rating float32 `json:"rating"`
		Note   string  `json:"note"`
	}

	PutMovieRecord struct {
		ID        int32  `param:"id"`
		RecordID  int32  `param:"recordId"`
		Date      string `json:"date"`
		MediaCode string `json:"mediaCode"`
	}
)

func NewMockedMovieHandler() routers.IRoute {
	return &mockedMovieHandler{}
}

func (h *mockedMovieHandler) Register(g *echo.Group) {
	g.GET("/movies", h.getMovies)
	g.GET("/movies/:id", h.getMovieDetail)

	g.POST("/movies/:id/impression", h.postMovieImpression)
	g.PUT("/movies/:id/impression", h.putMovieImpression)

	g.PUT("/movies/:id/records/:recordId", h.putMovieRecord)
}

func (h *mockedMovieHandler) getMovies(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("[Mock] accessed GET movies")

	return c.JSON(http.StatusOK, mock.MockedMovies)
}

func (h *mockedMovieHandler) getMovieDetail(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("[Mock] accessed GET movie detail")

	var req GetMovieDetailRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	logger.Info().Msg("successfully validated")

	result, ok := mock.MockedMovieDetailMapper[req.ID]
	if !ok {
		return c.String(http.StatusNotFound, fmt.Sprintf("movie(id=%d) is not found", req.ID))
	}

	return c.JSON(http.StatusOK, result)
}

func (h *mockedMovieHandler) postMovieImpression(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("[Mock] accessed POST movie impression")

	var req PostMovieImpression
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	logger.Info().Msg("successfully validated")

	return c.NoContent(http.StatusNoContent)
}

func (h *mockedMovieHandler) putMovieImpression(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("[Mock] accessed PUT movie impression")

	var req PutMovieImpression
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	logger.Info().Msg("successfully validated")

	return c.NoContent(http.StatusNoContent)
}

func (h *mockedMovieHandler) putMovieRecord(c echo.Context) error {
	logger := logger.GetLogger()
	logger.Info().Msg("[Mock] accessed PUT movie record")

	var req PutMovieRecord
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	logger.Info().Msg("successfully validated")

	return c.NoContent(http.StatusNoContent)
}
