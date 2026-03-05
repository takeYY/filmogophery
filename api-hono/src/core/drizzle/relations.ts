import { relations } from "drizzle-orm/relations";
import { movies, movieGenres, genres, series, users, refreshTokens, reviews, watchHistory, platforms, watchlist } from "./schema";

export const movieGenresRelations = relations(movieGenres, ({one}) => ({
	movie: one(movies, {
		fields: [movieGenres.movieId],
		references: [movies.id]
	}),
	genre: one(genres, {
		fields: [movieGenres.genreId],
		references: [genres.id]
	}),
}));

export const moviesRelations = relations(movies, ({one, many}) => ({
	movieGenres: many(movieGenres),
	series: one(series, {
		fields: [movies.seriesId],
		references: [series.id]
	}),
	reviews: many(reviews),
	watchHistories: many(watchHistory),
	watchlists: many(watchlist),
}));

export const genresRelations = relations(genres, ({many}) => ({
	movieGenres: many(movieGenres),
}));

export const seriesRelations = relations(series, ({many}) => ({
	movies: many(movies),
}));

export const refreshTokensRelations = relations(refreshTokens, ({one}) => ({
	user: one(users, {
		fields: [refreshTokens.userId],
		references: [users.id]
	}),
}));

export const usersRelations = relations(users, ({many}) => ({
	refreshTokens: many(refreshTokens),
	reviews: many(reviews),
	watchHistories: many(watchHistory),
	watchlists: many(watchlist),
}));

export const reviewsRelations = relations(reviews, ({one}) => ({
	user: one(users, {
		fields: [reviews.userId],
		references: [users.id]
	}),
	movie: one(movies, {
		fields: [reviews.movieId],
		references: [movies.id]
	}),
}));

export const watchHistoryRelations = relations(watchHistory, ({one}) => ({
	user: one(users, {
		fields: [watchHistory.userId],
		references: [users.id]
	}),
	movie: one(movies, {
		fields: [watchHistory.movieId],
		references: [movies.id]
	}),
	platform: one(platforms, {
		fields: [watchHistory.platformId],
		references: [platforms.id]
	}),
}));

export const platformsRelations = relations(platforms, ({many}) => ({
	watchHistories: many(watchHistory),
}));

export const watchlistRelations = relations(watchlist, ({one}) => ({
	user: one(users, {
		fields: [watchlist.userId],
		references: [users.id]
	}),
	movie: one(movies, {
		fields: [watchlist.movieId],
		references: [movies.id]
	}),
}));