import { dbConnections } from "@/core/db";
import { genres, movieGenres, movies, watchlist } from "@/core/drizzle/schema";
import { desc, eq, sql } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

/**
 * ユーザーのウォッチリストを映画・ジャンル情報付きで取得する
 */
export async function fetchWatchlistByUserId(
  userId: number,
  limit: number,
  offset: number,
  db: MySql2Database = dbConnections.readonly,
) {
  return db
    .select({
      id: watchlist.id,
      addedAt: watchlist.addedAt,
      priority: watchlist.priority,
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
    .from(watchlist)
    .innerJoin(movies, eq(watchlist.movieId, movies.id))
    .leftJoin(movieGenres, eq(movies.id, movieGenres.movieId))
    .leftJoin(genres, eq(movieGenres.genreId, genres.id))
    .where(eq(watchlist.userId, userId))
    .groupBy(watchlist.id)
    .orderBy(desc(watchlist.addedAt))
    .limit(limit)
    .offset(offset);
}

/**
 * 映画IDに一致する映画が存在するか確認する
 */
export async function fetchMovieById(
  movieId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const [movie] = await db
    .select({ id: movies.id })
    .from(movies)
    .where(eq(movies.id, movieId))
    .limit(1);
  return movie ?? null;
}

/**
 * ウォッチリストに登録する
 */
export async function insertWatchlist(
  userId: number,
  movieId: number,
  priority: number,
  db: MySql2Database = dbConnections.default,
) {
  await db.insert(watchlist).values({ userId, movieId, priority });
}

/**
 * ウォッチリストIDに一致するレコードを削除する
 * 削除件数を返す（0なら該当なし）
 */
export async function deleteWatchlistById(
  watchlistId: number,
  db: MySql2Database = dbConnections.default,
) {
  const [result] = await db
    .delete(watchlist)
    .where(eq(watchlist.id, watchlistId));
  return result.affectedRows;
}
