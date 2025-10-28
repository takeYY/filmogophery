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
			// --- Health --- //
			asV1Route(healthHandler.NewCheckHealthHandler), // check health
			// --- Movie --- //
			asV1Route(movieHandler.NewGetMovieDetailHandler), // get movie detail
			asV1Route(movieHandler.NewGetMoviesHandler),      // get movies
			// --- Review --- //
			asV1Route(reviewHandler.NewPostReviewHandler),        // create review
			asV1Route(reviewHandler.NewPutReviewHandler),         // update review
			asV1Route(reviewHandler.NewGetReviewHistoryHandler),  // get review history
			asV1Route(reviewHandler.NewPostReviewHistoryHandler), // create review history
			// --- Master --- //
			asV1Route(genreHandler.NewGetGenresHandler),       // get genres
			asV1Route(platformHandler.NewGetPlatformsHandler), // get platforms

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
