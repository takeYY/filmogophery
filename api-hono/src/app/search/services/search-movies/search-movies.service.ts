import { defaultRuntimeMinutes } from "@/core/definition";
import { redisService } from "@/core/services/redis/redis.service";
import { getMoviesByTitle } from "@/core/services/tmdb/tmdb.service";
import { Movie } from "@/core/types/movie";
import type { TmdbMovieResult } from "@/core/types/tmdb";
import { err, ok } from "neverthrow";
import { Logger } from "pino";
import {
  batchCreateMovies,
  getMoviesByTmdbIds,
} from "../../repositories/movies/movies.repository";

const CACHE_TTL_SECONDS = 24 * 60 * 60; // 24時間

export async function searchMovies(
  logger: Logger,
  title: string,
  limit: number,
  offset: number,
) {
  logger.info({ title, limit, offset }, "searchMovies called");

  // Redis 格納用のキャッシュキーを生成
  const cacheKey = newCacheKey(title, limit, offset);

  // Redis から情報を取得（あれば）
  const cached = await redisService.get<Movie[]>(cacheKey);
  if (cached !== null) {
    logger.debug("cache hit from redis");
    return ok(cached);
  }

  // Redis になければ TMDb API から映画情報を取得
  const tmdbResult = await getMoviesByTitle(title, offset);
  if (tmdbResult.isErr()) {
    return err(tmdbResult.error);
  }
  const tmdbMovies = tmdbResult.value.results.slice(0, limit);
  logger.debug("successfully search movies from tmdb");

  // 取得した映画の tmdbId リストで既存映画を取得
  const tmdbIds = tmdbMovies.map((m) => m.id);
  const existingMovies = await getMoviesByTmdbIds(tmdbIds);
  logger.debug("successfully get existing movies by tmdbId");

  // DB にない映画を新規登録
  const existingTmdbIdSet = new Set(existingMovies.map((m) => m.tmdbId));
  const newMovies = tmdbMovies
    .filter((m) => !existingTmdbIdSet.has(m.id))
    .map((m) => toMovieForCreation(m));

  if (newMovies.length > 0) {
    const newIds = await batchCreateMovies(newMovies);
    const newTmdbIds = newMovies.map((m) => m.tmdbId);
    const created = await getMoviesByTmdbIds(newTmdbIds);
    existingMovies.push(...created);
    logger.debug({ count: newIds.length }, "batch created new movies");
  }

  // TMDb の順序に合わせてレスポンス用に変換
  const movieMap = new Map(existingMovies.map((m) => [m.tmdbId, m]));
  const resultMovies: Movie[] = tmdbIds
    .map((tmdbId) => {
      const mv = movieMap.get(tmdbId);
      if (!mv) return null;
      const codes = mv.genreCodes?.split(",") ?? [];
      const names = mv.genreNames?.split(",") ?? [];
      return {
        id: mv.id,
        tmdbID: mv.tmdbId,
        title: mv.title,
        overview: mv.overview,
        releaseDate: mv.releaseDate,
        runtimeMinutes: mv.runtimeMinute,
        posterURL: mv.posterUrl ?? null,
        genres: codes.map((code, i) => ({ code, name: names[i] ?? "" })),
      } satisfies Movie;
    })
    .filter((m): m is Movie => m !== null);

  // Redis にキャッシュ（24時間）
  await redisService
    .set(cacheKey, resultMovies, CACHE_TTL_SECONDS)
    .catch((e) => {
      logger.warn({ err: e }, "failed to cache movies in redis");
    });

  return ok(resultMovies);
}

function newCacheKey(title: string, limit: number, offset: number): string {
  return `movies:search:${title.trim().toLowerCase()}:limit:${limit}:offset:${offset}`;
}

function toMovieForCreation(m: TmdbMovieResult) {
  return {
    tmdbId: m.id,
    title: m.title,
    overview: m.overview,
    releaseDate: m.release_date || "1970-01-01",
    runtimeMinutes: defaultRuntimeMinutes,
    posterUrl: m.poster_path,
    genreIds: m.genre_ids,
  };
}
