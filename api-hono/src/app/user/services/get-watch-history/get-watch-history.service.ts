import { ok } from "neverthrow";
import { Logger } from "pino";
import { fetchWatchHistoryByUserId } from "../../repositories/watch-history/watch-history.repository";

export type WatchHistoryItem = {
  id: number;
  watchedAt: string | null;
  platform: {
    id: number;
    code: string;
    name: string;
  };
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

export async function getWatchHistory(
  logger: Logger,
  userId: number,
  limit: number,
  offset: number,
) {
  logger.info({ userId, limit, offset }, "getWatchHistory called");

  const rows = await fetchWatchHistoryByUserId(userId, limit, offset);

  const result: WatchHistoryItem[] = rows.map((row) => {
    const codes = row.genreCodes?.split(",") || [];
    const names = row.genreNames?.split(",") || [];
    const genres = codes.map((code, i) => ({ code, name: names[i] ?? "" }));

    return {
      id: row.id,
      watchedAt: row.watchedDate ?? null,
      platform: {
        id: row.platformId,
        code: row.platformCode,
        name: row.platformName,
      },
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

  logger.info({ count: result.length }, "getWatchHistory completed");
  return ok(result);
}
