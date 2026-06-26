import { MovieIsNotFound } from "@/core/errors";
import { err, ok, Result } from "neverthrow";
import { Logger } from "pino";
import { fetchMovieById } from "../../repositories/movies/movies.repository";
import { fetchWatchHistoryByMovieId } from "../../repositories/watch-history/watch-history.repository";

export type WatchHistoryItem = {
  id: number;
  platform: {
    id: number;
    code: string;
    name: string;
  };
  watchedAt: string | null;
};

export async function getMovieWatchHistory(
  logger: Logger,
  userId: number,
  movieId: number,
): Promise<Result<WatchHistoryItem[], MovieIsNotFound>> {
  logger.info({ userId, movieId }, "getMovieWatchHistory called");

  // 映画の存在確認
  const movie = await fetchMovieById(movieId);
  if (!movie) {
    logger.info({ movieId }, "movie not found");
    return err(new MovieIsNotFound(movieId));
  }

  const rows = await fetchWatchHistoryByMovieId(userId, movieId);

  const result: WatchHistoryItem[] = rows.map((row) => ({
    id: row.id,
    platform: {
      id: row.platformId,
      code: row.platformCode,
      name: row.platformName,
    },
    watchedAt: row.watchedDate ?? null,
  }));

  logger.info({ count: result.length }, "getMovieWatchHistory completed");
  return ok(result);
}
