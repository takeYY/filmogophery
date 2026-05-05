import { dbConnections } from "@/core/db";
import { platforms, watchHistory } from "@/core/drizzle/schema";
import { and, eq } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

/**
 * ユーザーの映画視聴履歴をプラットフォーム情報付きで取得する
 */
export async function fetchWatchHistoryByMovieId(
  userId: number,
  movieId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  return db
    .select({
      id: watchHistory.id,
      watchedDate: watchHistory.watchedDate,
      platformId: platforms.id,
      platformCode: platforms.code,
      platformName: platforms.name,
    })
    .from(watchHistory)
    .innerJoin(platforms, eq(watchHistory.platformId, platforms.id))
    .where(
      and(eq(watchHistory.userId, userId), eq(watchHistory.movieId, movieId)),
    );
}
