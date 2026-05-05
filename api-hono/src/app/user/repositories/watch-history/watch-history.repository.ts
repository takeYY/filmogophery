import { dbConnections } from "@/core/db";
import {
  genres,
  movieGenres,
  movies,
  platforms,
  watchHistory,
} from "@/core/drizzle/schema";
import { and, desc, eq, sql } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

/**
 * ユーザーの視聴履歴を映画・プラットフォーム・ジャンル情報付きで取得する
 */
export async function fetchWatchHistoryByUserId(
  userId: number,
  limit: number,
  offset: number,
  db: MySql2Database = dbConnections.readonly,
) {
  // 視聴履歴と映画・プラットフォームを結合して取得
  const rows = await db
    .select({
      id: watchHistory.id,
      watchedDate: watchHistory.watchedDate,
      platformId: platforms.id,
      platformCode: platforms.code,
      platformName: platforms.name,
      movieId: movies.id,
      movieTitle: movies.title,
      movieOverview: movies.overview,
      movieReleaseDate: movies.releaseDate,
      movieRuntimeMinutes: movies.runtimeMinutes,
      moviePosterUrl: movies.posterUrl,
      movieTmdbId: movies.tmdbId,
      genreCodes: sql<string>`GROUP_CONCAT(DISTINCT ${genres.code})`.as(
        "genre_codes",
      ),
      genreNames: sql<string>`GROUP_CONCAT(DISTINCT ${genres.name})`.as(
        "genre_names",
      ),
    })
    .from(watchHistory)
    .innerJoin(platforms, eq(watchHistory.platformId, platforms.id))
    .innerJoin(movies, eq(watchHistory.movieId, movies.id))
    .leftJoin(movieGenres, eq(movies.id, movieGenres.movieId))
    .leftJoin(genres, eq(movieGenres.genreId, genres.id))
    .where(and(eq(watchHistory.userId, userId)))
    .groupBy(watchHistory.id)
    .orderBy(desc(watchHistory.watchedDate))
    .limit(limit)
    .offset(offset);

  return rows;
}
