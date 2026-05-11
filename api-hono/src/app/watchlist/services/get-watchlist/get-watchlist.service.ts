import { ok } from "neverthrow";
import { Logger } from "pino";
import { fetchWatchlistByUserId } from "../../repositories/watchlist.repository";

export type WatchlistItem = {
  id: number;
  addedAt: string | null;
  priority: number | null;
  movie: {
    id: number;
    title: string;
    overview: string;
    releaseDate: string;
    runtimeMinutes: number;
    posterURL: string | null;
    tmdbID: number;
    genres: { code: string; name: string }[];
  };
};

export async function getWatchlist(
  logger: Logger,
  userId: number,
  limit: number,
  offset: number,
) {
  logger.info({ userId, limit, offset }, "getWatchlist called");

  const rows = await fetchWatchlistByUserId(userId, limit, offset);

  const result: WatchlistItem[] = rows.map((row) => {
    const codes = row.genreCodes?.split(",") || [];
    const names = row.genreNames?.split(",") || [];
    const genres = codes.map((code, i) => ({ code, name: names[i] ?? "" }));

    return {
      id: row.id,
      addedAt: row.addedAt ?? null,
      priority: row.priority ?? null,
      movie: {
        id: row.movieId,
        title: row.movieTitle,
        overview: row.movieOverview,
        releaseDate: row.movieReleaseDate,
        runtimeMinutes: row.movieRuntimeMinutes,
        posterURL: row.moviePosterUrl ?? null,
        tmdbID: row.movieTmdbId,
        genres,
      },
    };
  });

  logger.info({ count: result.length }, "getWatchlist completed");
  return ok(result);
}
