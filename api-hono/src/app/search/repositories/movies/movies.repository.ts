import { dbConnections } from "@/core/db";
import { genres, movieGenres, movies } from "@/core/drizzle/schema";
import { inArray, sql } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function getMoviesByTmdbIds(
  tmdbIds: number[],
  db: MySql2Database = dbConnections.readonly,
) {
  return db
    .select({
      id: movies.id,
      tmdbId: movies.tmdbId,
      title: movies.title,
      overview: movies.overview,
      releaseDate: movies.releaseDate,
      runtimeMinute: movies.runtimeMinutes,
      posterUrl: movies.posterUrl,
      genreCodes: sql<string>`GROUP_CONCAT(DISTINCT ${genres.code})`.as(
        "genre_codes",
      ),
      genreNames: sql<string>`GROUP_CONCAT(DISTINCT ${genres.name})`.as(
        "genre_names",
      ),
    })
    .from(movies)
    .leftJoin(movieGenres, sql`${movies.id} = ${movieGenres.movieId}`)
    .leftJoin(genres, sql`${movieGenres.genreId} = ${genres.id}`)
    .where(inArray(movies.tmdbId, tmdbIds))
    .groupBy(movies.id);
}

export async function batchCreateMovies(
  newMovies: Array<{
    tmdbId: number;
    title: string;
    overview: string;
    releaseDate: string;
    runtimeMinutes: number;
    posterUrl: string | null;
    genreIds: number[];
  }>,
  db: MySql2Database = dbConnections.default,
) {
  return db.transaction(async (tx) => {
    const created: number[] = [];

    for (const mv of newMovies) {
      const [result] = await tx.insert(movies).values({
        id: mv.tmdbId,
        tmdbId: mv.tmdbId,
        title: mv.title,
        overview: mv.overview,
        releaseDate: mv.releaseDate,
        runtimeMinutes: mv.runtimeMinutes,
        posterUrl: mv.posterUrl,
      });

      const insertId = result.insertId;
      created.push(insertId);

      if (mv.genreIds.length > 0) {
        await tx
          .insert(movieGenres)
          .values(
            mv.genreIds.map((genreId) => ({ movieId: insertId, genreId })),
          );
      }
    }

    return created;
  });
}
