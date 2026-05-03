import { users } from "@/core/drizzle/schema";
import { MovieIsNotFound, TmdbMovieIsNotFound } from "@/core/errors";
import { getTmdbMovieDetailById } from "@/core/services/tmdb/tmdb.service";
import { Movie, MovieDetail, Review, Series } from "@/core/types/movie";
import type { TmdbMovieDetailResponse } from "@/core/types/tmdb";
import { err, Ok, ok } from "neverthrow";
import { Logger } from "pino";
import {
  fetchMovieById,
  getReviewedMoviesByUser,
} from "../../repositories/movies/movies.repository";
import { fetchReviewByMovieId } from "../../repositories/reviews/reviews.repository";
import { fetchSeriesByMovieId } from "../../repositories/series/series.repository";

export async function getMovies(
  logger: Logger,
  userId: number,
  genre: string | undefined,
  limit: number,
  offset: number,
): Promise<Ok<Movie[], never>> {
  logger.info({ userId, genre, limit, offset }, "getMovies called");
  const result = await getReviewedMoviesByUser(userId, genre, limit, offset);
  if (result.length == 0) {
    logger.info("movie is not found");
    return ok([]);
  }

  const movies: Movie[] = result.map((movie) => {
    const codes = movie.genreCodes?.split(",") || [];
    const names = movie.genreNames?.split(",") || [];
    const genres = codes.map((code, i) => ({ code, name: names[i] }));

    return {
      id: movie.id,
      title: movie.title,
      overview: movie.overview,
      releaseDate: movie.releaseDate,
      runtimeMinute: movie.runtimeMinute,
      posterUrl: movie.posterUrl,
      tmdbId: movie.tmdbId,
      genres,
    };
  });

  return ok(movies);
}

export async function getMovieById(
  logger: Logger,
  operator: typeof users.$inferSelect,
  id: number,
) {
  logger.info({ id }, "getMovieById called");

  const movie = await fetchMovieById(id);
  if (movie === undefined) {
    logger.info({ id }, "movie is not found");
    return err(new MovieIsNotFound());
  }

  const codes = movie.genreCodes?.split(",") || [];
  const names = movie.genreNames?.split(",") || [];
  const genres = codes.map((code, i) => ({ code, name: names[i] }));

  const tmdbMovie = await getTmdbMovieDetailById(movie.tmdbId);
  if (tmdbMovie.isErr()) {
    return err(new TmdbMovieIsNotFound());
  }
  const tmdb = (await tmdbMovie.value.json()) as TmdbMovieDetailResponse;

  const review = await fetchReviewByMovieId(operator.id, movie.id);
  let reviewResponse: Review | null;
  if (review === undefined) {
    reviewResponse = null;
  } else {
    reviewResponse = {
      id: review.id,
      rating: review.rating,
      comment: review.comment,
      createdAt: review.createdAt,
      updatedAt: review.updatedAt,
    };
  }

  const voteAverage = Math.round((tmdb.vote_average / 2) * 10) / 10;

  const series = await fetchSeriesByMovieId(movie.id);
  let seriesResponse: Series | null;
  if (series === undefined) {
    seriesResponse = null;
  } else {
    seriesResponse = {
      name: series.name,
      posterUrl: series.posterUrl,
    };
  }

  const result: MovieDetail = {
    ...movie,
    voteAverage,
    voteCount: tmdb.vote_count,
    series: seriesResponse,
    review: reviewResponse,
    genres,
  };

  return ok(result);
}
