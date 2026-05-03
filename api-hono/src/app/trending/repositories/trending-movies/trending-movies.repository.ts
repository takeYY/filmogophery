import { dbConnections } from "@/core/db";
import { movies, reviews } from "@/core/drizzle/schema";
import { and, eq, inArray } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

/**
 * tmdbId のリストに一致する映画を取得する
 */
export async function fetchMoviesByTmdbIds(
  tmdbIds: number[],
  db: MySql2Database = dbConnections.readonly,
) {
  if (tmdbIds.length === 0) return [];

  return db
    .select({
      id: movies.id,
      tmdbId: movies.tmdbId,
      title: movies.title,
      posterUrl: movies.posterUrl,
    })
    .from(movies)
    .where(inArray(movies.tmdbId, tmdbIds));
}

/**
 * 映画を一括登録する（既存レコードは無視）
 */
export async function batchInsertMovies(
  values: {
    tmdbId: number;
    title: string;
    overview: string;
    releaseDate: string;
    runtimeMinutes: number;
    posterUrl: string | null;
  }[],
  db: MySql2Database = dbConnections.default,
) {
  if (values.length === 0) return;

  // INSERT IGNORE で重複をスキップ
  await db
    .insert(movies)
    .ignore()
    .values(
      values.map((v) => ({
        tmdbId: v.tmdbId,
        title: v.title,
        overview: v.overview,
        releaseDate: v.releaseDate,
        runtimeMinutes: v.runtimeMinutes,
        posterUrl: v.posterUrl,
      })),
    );
}

/**
 * 映画IDリストのうちユーザーがレビュー済みの映画IDセットを取得する
 */
export async function fetchReviewedMovieIds(
  userId: number,
  movieIds: number[],
  db: MySql2Database = dbConnections.readonly,
): Promise<Set<number>> {
  if (movieIds.length === 0) return new Set();

  const rows = await db
    .select({ movieId: reviews.movieId })
    .from(reviews)
    .where(and(eq(reviews.userId, userId), inArray(reviews.movieId, movieIds)));

  return new Set(rows.map((r) => r.movieId));
}
