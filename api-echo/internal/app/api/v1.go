package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	authHandler "filmogophery/internal/app/features/auth/handlers"
	genreHandler "filmogophery/internal/app/features/genre/handlers"
	healthHandler "filmogophery/internal/app/features/health/handlers"
	movieHandler "filmogophery/internal/app/features/movie/handlers"
	platformHandler "filmogophery/internal/app/features/platform/handlers"
	reviewHandler "filmogophery/internal/app/features/review/handlers"
	searchHandler "filmogophery/internal/app/features/search/handlers"
	trendingHandler "filmogophery/internal/app/features/trending/handlers"
	userHandler "filmogophery/internal/app/features/user/handlers"
	watchlistHandler "filmogophery/internal/app/features/watchlist/handlers"
	"filmogophery/internal/app/routers"
)

func RegisterV1Routes() fx.Option {
	return fx.Module(
		"v1-route",
		fx.Provide(
			fx.Private,
			// --- Health --- //
			asV1Route(healthHandler.NewCheckHealthHandler), // check health
			// --- User --- //
			asV1Route(userHandler.NewCreateUserHandler), // create user
			// --- Auth --- //
			asV1Route(authHandler.NewLoginHandler), // login
			// --- Movie --- //
			asV1Route(movieHandler.NewGetMovieDetailHandler),       // get movie detail
			asV1Route(movieHandler.NewGetMoviesHandler),            // get movies
			asV1Route(movieHandler.NewGetMovieWatchHistoryHandler), // get movie watch history
			// --- Review --- //
			asV1Route(reviewHandler.NewPostReviewHandler),        // create review
			asV1Route(reviewHandler.NewPutReviewHandler),         // update review
			asV1Route(reviewHandler.NewGetReviewHistoryHandler),  // get review history
			asV1Route(reviewHandler.NewPostReviewHistoryHandler), // create review history
			// --- Watchlist --- //
			asV1Route(watchlistHandler.NewGetWatchlistHandler), // get watchlist
			// --- Trending --- //
			asV1Route(trendingHandler.NewGetTrendingMoviesHandler), // get trending movies
			// --- Search --- //
			asV1Route(searchHandler.NewSearchMoviesHandler), // search movies
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
