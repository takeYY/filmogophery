import { dbConnections } from "@/core/db";
import { genres, movieGenres, movies, reviews } from "@/core/drizzle/schema";
import { and, eq, sql } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function getMoviesByGenre(
  genre: string | undefined,
  limit: number,
  offset: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const query = db
    .select({
      id: movies.id,
      title: movies.title,
      overview: movies.overview,
      releaseDate: movies.releaseDate,
      runtimeMinute: movies.runtimeMinutes,
      posterUrl: movies.posterUrl,
      tmdbId: movies.tmdbId,
      genreCodes: sql<string>`GROUP_CONCAT(DISTINCT ${genres.code})`.as(
        "genre_codes",
      ),
      genreNames: sql<string>`GROUP_CONCAT(DISTINCT ${genres.name})`.as(
        "genre_names",
      ),
    })
    .from(movies)
    .leftJoin(movieGenres, eq(movies.id, movieGenres.movieId))
    .leftJoin(genres, eq(movieGenres.genreId, genres.id))
    .groupBy(movies.id);

  if (genre) {
    return await query
      .having(sql`FIND_IN_SET(${genre}, GROUP_CONCAT(DISTINCT ${genres.code}))`)
      .limit(limit)
      .offset(offset);
  }

  return await query.limit(limit).offset(offset);
}

export async function getReviewedMoviesByUser(
  userId: number,
  genre: string | undefined,
  limit: number,
  offset: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const query = db
    .select({
      id: movies.id,
      title: movies.title,
      overview: movies.overview,
      releaseDate: movies.releaseDate,
      runtimeMinute: movies.runtimeMinutes,
      posterUrl: movies.posterUrl,
      tmdbId: movies.tmdbId,
      genreCodes: sql<string>`GROUP_CONCAT(DISTINCT ${genres.code})`.as(
        "genre_codes",
      ),
      genreNames: sql<string>`GROUP_CONCAT(DISTINCT ${genres.name})`.as(
        "genre_names",
      ),
    })
    .from(movies)
    .innerJoin(
      reviews,
      and(eq(reviews.movieId, movies.id), eq(reviews.userId, userId)),
    )
    .leftJoin(movieGenres, eq(movies.id, movieGenres.movieId))
    .leftJoin(genres, eq(movieGenres.genreId, genres.id))
    .groupBy(movies.id)
    .orderBy(sql`MAX(${reviews.createdAt}) DESC`);

  if (genre) {
    return await query
      .having(sql`FIND_IN_SET(${genre}, GROUP_CONCAT(DISTINCT ${genres.code}))`)
      .limit(limit)
      .offset(offset);
  }

  return await query.limit(limit).offset(offset);
}

export async function fetchMovieById(
  id: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const [movie] = await db
    .select({
      id: movies.id,
      title: movies.title,
      overview: movies.overview,
      releaseDate: movies.releaseDate,
      runtimeMinute: movies.runtimeMinutes,
      posterUrl: movies.posterUrl,
      tmdbId: movies.tmdbId,
      genreCodes: sql<string>`GROUP_CONCAT(DISTINCT ${genres.code})`.as(
        "genre_codes",
      ),
      genreNames: sql<string>`GROUP_CONCAT(DISTINCT ${genres.name})`.as(
        "genre_names",
      ),
    })
    .from(movies)
    .leftJoin(movieGenres, eq(movies.id, movieGenres.movieId))
    .leftJoin(genres, eq(movieGenres.genreId, genres.id))
    .where(eq(movies.id, id))
    .groupBy(movies.id);

  return movie;
}

export async function updateRuntimeMinutes(
  id: number,
  runtimeMinutes: number,
  db: MySql2Database = dbConnections.default,
) {
  const [result] = await db
    .update(movies)
    .set({ runtimeMinutes })
    .where(eq(movies.id, id));

  return result.affectedRows > 0;
}
