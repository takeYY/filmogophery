import { dbConnections } from "@/core/db";
import { genres, movieGenres, movies } from "@/core/drizzle/schema";
import { eq, sql } from "drizzle-orm";
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
