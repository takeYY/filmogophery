import { dbConnections } from "@/core/db";
import { movies, series } from "@/core/drizzle/schema";
import { eq } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function fetchSeriesByMovieId(
  movieId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  return (
    await db
      .select({
        name: series.name,
        posterUrl: series.posterUrl,
      })
      .from(movies)
      .leftJoin(series, eq(series.id, movies.seriesId))
      .where(eq(movies.id, movieId))
      .limit(1)
  ).at(0);
}
