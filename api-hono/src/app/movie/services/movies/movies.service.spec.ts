import * as movieRepo from "@/app/movie/repositories/movies/movies.repository";
import { getMovies } from "@/app/movie/services/movies/movies.service";
import { afterEach, describe, expect, mock, spyOn, test } from "bun:test";
import pino from "pino";

const logger = pino({ level: "silent" });
const userId = 1;

describe("getMovies", () => {
  afterEach(() => {
    mock.restore();
  });

  test("レビュー済み映画が見つからない場合は空配列を返す", async () => {
    spyOn(movieRepo, "getReviewedMoviesByUser").mockResolvedValue([]);

    const result = await getMovies(logger, userId, undefined, 12, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toBeArrayOfSize(0);
  });

  test("レビュー済み映画が見つかった場合はジャンルをオブジェクトで返す", async () => {
    spyOn(movieRepo, "getReviewedMoviesByUser").mockResolvedValue([
      {
        id: 1,
        title: "テスト映画",
        overview: "テスト映画概要",
        releaseDate: "2026-03-03",
        runtimeMinute: 314,
        posterUrl: null,
        tmdbId: 2,
        genreCodes: "action,sf",
        genreNames: "アクション,SF",
      },
    ]);

    const result = await getMovies(logger, userId, undefined, 12, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toEqual([
      {
        id: 1,
        title: "テスト映画",
        overview: "テスト映画概要",
        releaseDate: "2026-03-03",
        runtimeMinutes: 314,
        posterURL: null,
        tmdbID: 2,
        genres: [
          { code: "action", name: "アクション" },
          { code: "sf", name: "SF" },
        ],
      },
    ]);
  });

  test("ジャンル絞り込みありでレビュー済み映画を返す", async () => {
    spyOn(movieRepo, "getReviewedMoviesByUser").mockResolvedValue([
      {
        id: 1,
        title: "アクション映画",
        overview: "概要",
        releaseDate: "2026-03-03",
        runtimeMinute: 120,
        posterUrl: null,
        tmdbId: 3,
        genreCodes: "action",
        genreNames: "アクション",
      },
    ]);

    const result = await getMovies(logger, userId, "action", 12, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toHaveLength(1);
    expect(data[0].genres).toEqual([{ code: "action", name: "アクション" }]);
  });
});
