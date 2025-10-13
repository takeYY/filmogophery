package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	genreHandler "filmogophery/internal/app/features/genre/handlers"
	healthHandler "filmogophery/internal/app/features/health/handlers"
	movieHandler "filmogophery/internal/app/features/movie/handlers"
	platformHandler "filmogophery/internal/app/features/platform/handlers"
	reviewHandler "filmogophery/internal/app/features/review/handlers"
	"filmogophery/internal/app/routers"
)

func RegisterV1Routes() fx.Option {
	return fx.Module(
		"v1-route",
		fx.Provide(
			fx.Private,
			asV1Route(healthHandler.NewCheckHealthHandler),
			asV1Route(movieHandler.NewGetMovieDetailHandler),
			asV1Route(movieHandler.NewGetMoviesHandler),
			asV1Route(reviewHandler.NewGetReviewHistoryHandler),
			asV1Route(genreHandler.NewGetGenresHandler),
			asV1Route(platformHandler.NewGetPlatformsHandler),

			fx.Annotate(
				getRouters,
				fx.ParamTags(`group:"v1-routers"`),
			),
		),
		fx.Invoke(func(e *echo.Echo, hs []routers.IRoute, m ...echo.MiddlewareFunc) {
			g := e.Group("/v1", m...)
			for _, h := range hs {
				h.Register(g)
			}
		}),
	)
}

func asV1Route(h any) any {
	return fx.Annotate(
		h,
		fx.As(new(routers.IRoute)),
		fx.ResultTags(`group:"v1-routers"`),
	)
}

func getRouters(h []routers.IRoute) []routers.IRoute {
	return h
}
