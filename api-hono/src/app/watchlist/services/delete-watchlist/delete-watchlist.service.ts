import { err, ok, Result } from "neverthrow";
import { Logger } from "pino";
import { deleteWatchlistById } from "../../repositories/watchlist.repository";

export class WatchlistNotFoundError extends Error {}

export async function deleteWatchlist(
  logger: Logger,
  watchlistId: number,
): Promise<Result<void, WatchlistNotFoundError>> {
  logger.info({ watchlistId }, "deleteWatchlist called");

  const affected = await deleteWatchlistById(watchlistId);
  if (affected === 0) {
    logger.info({ watchlistId }, "watchlist not found");
    return err(new WatchlistNotFoundError());
  }

  logger.info({ watchlistId }, "deleteWatchlist completed");
  return ok(undefined);
}
