import { err, ok, Result } from "neverthrow";
import { Logger } from "pino";
import {
  fetchMovieById,
  insertWatchlist,
} from "../../repositories/watchlist.repository";

export class MovieNotFoundError extends Error {}

export async function addWatchlist(
  logger: Logger,
  userId: number,
  movieId: number,
  priority: number,
): Promise<Result<void, MovieNotFoundError>> {
  logger.info({ userId, movieId, priority }, "addWatchlist called");

  // 映画の存在確認
  const movie = await fetchMovieById(movieId);
  if (!movie) {
    logger.info({ movieId }, "movie not found");
    return err(new MovieNotFoundError());
  }

  await insertWatchlist(userId, movieId, priority);

  logger.info({ userId, movieId }, "addWatchlist completed");
  return ok(undefined);
}
