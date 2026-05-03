import { getTrendingMovies as fetchTmdbTrending } from "@/core/services/tmdb/tmdb.service";
import { TrendingMovie } from "@/core/types/movie";
import { TmdbTrendingMovieResult } from "@/core/types/tmdb";
import { err, ok, Result } from "neverthrow";
import { Logger } from "pino";
import {
  batchInsertMovies,
  fetchMoviesByTmdbIds,
  fetchReviewedMovieIds,
} from "../../repositories/trending-movies/trending-movies.repository";

export class TrendingMoviesError extends Error {}

export async function getTrendingMovies(
  logger: Logger,
  userId: number,
): Promise<Result<TrendingMovie[], TrendingMoviesError>> {
  logger.info({ userId }, "getTrendingMovies called");

  // TMDBからトレンド映画を取得
  const tmdbResult = await fetchTmdbTrending();
  if (tmdbResult.isErr()) {
    logger.error(
      { err: tmdbResult.error },
      "failed to fetch trending from tmdb",
    );
    return err(new TrendingMoviesError("failed to fetch trending movies"));
  }

  const tmdbMovies = tmdbResult.value.results;
  if (tmdbMovies.length === 0) {
    return ok([]);
  }

  const tmdbIds = tmdbMovies.map((m) => m.id);

  // DBに存在する映画を取得
  const existingMovies = await fetchMoviesByTmdbIds(tmdbIds);
  const existingTmdbIdSet = new Set(existingMovies.map((m) => m.tmdbId));

  // DBに存在しない映画を登録
  const newMovies = tmdbMovies
    .filter((m: TmdbTrendingMovieResult) => !existingTmdbIdSet.has(m.id))
    .map((m: TmdbTrendingMovieResult) => ({
      tmdbId: m.id,
      title: m.title,
      overview: m.overview,
      releaseDate: m.release_date || "1970-01-01",
      runtimeMinutes: 1, // TMDBトレンドAPIでは上映時間が取れないため仮値
      posterUrl: m.poster_path,
    }));

  if (newMovies.length > 0) {
    await batchInsertMovies(newMovies);
    logger.info({ count: newMovies.length }, "inserted new trending movies");
  }

  // 登録後に全tmdbIdで再取得（新規登録分のIDを得るため）
  const allMovies = await fetchMoviesByTmdbIds(tmdbIds);
  const movieByTmdbId = new Map(allMovies.map((m) => [m.tmdbId, m]));

  // レビュー済みフラグを一括取得
  const allMovieIds = allMovies.map((m) => m.id);
  const reviewedIds = await fetchReviewedMovieIds(userId, allMovieIds);

  // TMDBの順序を維持してレスポンスを構築
  const result: TrendingMovie[] = [];
  for (const tmdbMovie of tmdbMovies) {
    const movie = movieByTmdbId.get(tmdbMovie.id);
    if (!movie) continue;
    result.push({
      id: movie.id,
      title: movie.title,
      posterUrl: movie.posterUrl ?? null,
      tmdbId: movie.tmdbId,
      hasReview: reviewedIds.has(movie.id),
    });
  }

  logger.info({ count: result.length }, "getTrendingMovies completed");
  return ok(result);
}
