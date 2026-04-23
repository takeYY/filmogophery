import * as moviesRepo from "@/app/search/repositories/movies/movies.repository";
import { searchMovies } from "@/app/search/services/search-movies/search-movies.service";
import * as redisServiceModule from "@/core/services/redis/redis.service";
import * as tmdbService from "@/core/services/tmdb/tmdb.service";
import { afterEach, describe, expect, mock, spyOn, test } from "bun:test";
import { err, ok } from "neverthrow";
import pino from "pino";

const logger = pino({ level: "silent" });

describe("searchMovies", () => {
  afterEach(() => {
    mock.restore();
  });

  test("キャッシュがある場合はキャッシュの値を返す", async () => {
    const cachedMovies = [
      {
        id: 1,
        tmdbId: 100,
        title: "キャッシュ映画",
        overview: "概要",
        releaseDate: "2024-01-01",
        runtimeMinute: 120,
        posterUrl: null,
        genres: [{ code: "action", name: "アクション" }],
      },
    ];

    spyOn(redisServiceModule.redisService, "get").mockResolvedValue(
      cachedMovies,
    );

    const result = await searchMovies(logger, "キャッシュ映画", 10, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toEqual(cachedMovies);
  });

  test("TMDb API がエラーを返した場合はエラーを返す", async () => {
    spyOn(redisServiceModule.redisService, "get").mockResolvedValue(null);
    spyOn(tmdbService, "getMoviesByTitle").mockResolvedValue(
      err(new Error("TMDb error")),
    );

    const result = await searchMovies(logger, "エラー映画", 10, 0);

    expect(result.isErr()).toBe(true);
  });

  test("DBに存在する映画はそのまま返す", async () => {
    spyOn(redisServiceModule.redisService, "get").mockResolvedValue(null);
    spyOn(tmdbService, "getMoviesByTitle").mockResolvedValue(
      ok({
        results: [
          {
            id: 100,
            title: "既存映画",
            overview: "概要",
            release_date: "2024-01-01",
            poster_path: null,
            genre_ids: [1],
          },
        ],
      } as any),
    );
    spyOn(moviesRepo, "getMoviesByTmdbIds").mockResolvedValue([
      {
        id: 1,
        tmdbId: 100,
        title: "既存映画",
        overview: "概要",
        releaseDate: "2024-01-01",
        runtimeMinute: 120,
        posterUrl: null,
        genreCodes: "action",
        genreNames: "アクション",
      },
    ]);
    spyOn(redisServiceModule.redisService, "set").mockResolvedValue(undefined);

    const result = await searchMovies(logger, "既存映画", 10, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toHaveLength(1);
    expect(data[0].tmdbId).toBe(100);
    expect(data[0].title).toBe("既存映画");
    expect(data[0].genres).toEqual([{ code: "action", name: "アクション" }]);
  });

  test("DBに存在しない映画は新規登録して返す", async () => {
    spyOn(redisServiceModule.redisService, "get").mockResolvedValue(null);
    spyOn(tmdbService, "getMoviesByTitle").mockResolvedValue(
      ok({
        results: [
          {
            id: 200,
            title: "新規映画",
            overview: "新規概要",
            release_date: "2024-06-01",
            poster_path: "/poster.jpg",
            genre_ids: [2],
          },
        ],
      } as any),
    );

    // 最初の呼び出し（既存チェック）は空、2回目（新規作成後）は登録済みを返す
    spyOn(moviesRepo, "getMoviesByTmdbIds")
      .mockResolvedValueOnce([])
      .mockResolvedValueOnce([
        {
          id: 10,
          tmdbId: 200,
          title: "新規映画",
          overview: "新規概要",
          releaseDate: "2024-06-01",
          runtimeMinute: 1,
          posterUrl: "/poster.jpg",
          genreCodes: "",
          genreNames: "",
        },
      ]);

    spyOn(moviesRepo, "batchCreateMovies").mockResolvedValue([10]);
    spyOn(redisServiceModule.redisService, "set").mockResolvedValue(undefined);

    const result = await searchMovies(logger, "新規映画", 10, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toHaveLength(1);
    expect(data[0].tmdbId).toBe(200);
    expect(data[0].title).toBe("新規映画");
    expect(moviesRepo.batchCreateMovies).toHaveBeenCalledTimes(1);
  });

  test("TMDb の結果が空の場合は空配列を返す", async () => {
    spyOn(redisServiceModule.redisService, "get").mockResolvedValue(null);
    spyOn(tmdbService, "getMoviesByTitle").mockResolvedValue(
      ok({ results: [] } as any),
    );
    spyOn(moviesRepo, "getMoviesByTmdbIds").mockResolvedValue([]);
    spyOn(redisServiceModule.redisService, "set").mockResolvedValue(undefined);

    const result = await searchMovies(logger, "存在しない映画", 10, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toBeArrayOfSize(0);
  });

  test("ジャンルが複数ある場合はオブジェクト配列で返す", async () => {
    spyOn(redisServiceModule.redisService, "get").mockResolvedValue(null);
    spyOn(tmdbService, "getMoviesByTitle").mockResolvedValue(
      ok({
        results: [
          {
            id: 300,
            title: "マルチジャンル映画",
            overview: "概要",
            release_date: "2024-01-01",
            poster_path: null,
            genre_ids: [1, 2],
          },
        ],
      } as any),
    );
    spyOn(moviesRepo, "getMoviesByTmdbIds").mockResolvedValue([
      {
        id: 5,
        tmdbId: 300,
        title: "マルチジャンル映画",
        overview: "概要",
        releaseDate: "2024-01-01",
        runtimeMinute: 90,
        posterUrl: null,
        genreCodes: "action,drama",
        genreNames: "アクション,ドラマ",
      },
    ]);
    spyOn(redisServiceModule.redisService, "set").mockResolvedValue(undefined);

    const result = await searchMovies(logger, "マルチジャンル映画", 10, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data[0].genres).toEqual([
      { code: "action", name: "アクション" },
      { code: "drama", name: "ドラマ" },
    ]);
  });
});
