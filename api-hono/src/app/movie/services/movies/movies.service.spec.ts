import * as movieRepo from "@/app/movie/repositories/movies/movies.repository";
import { getMovies } from "@/app/movie/services/movies/movies.service";
import { afterEach, describe, expect, mock, spyOn, test } from "bun:test";
import pino from "pino";

const logger = pino({ level: "silent" });

describe("getMovies", () => {
  afterEach(() => {
    mock.restore();
  });

  test("映画が見つからない場合は空配列を返す", async () => {
    spyOn(movieRepo, "getMoviesByGenre").mockResolvedValue([]);

    const result = await getMovies(logger, undefined, 12, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toBeArrayOfSize(0);
  });

  test("映画が見つかったあ場合はジャンルをオブジェクトで返す", async () => {
    spyOn(movieRepo, "getMoviesByGenre").mockResolvedValue([
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

    const result = await getMovies(logger, undefined, 12, 0);
    const data = result._unsafeUnwrap();

    expect(result.isOk()).toBe(true);
    expect(data).toEqual([
      {
        id: 1,
        title: "テスト映画",
        overview: "テスト映画概要",
        releaseDate: "2026-03-03",
        runtimeMinute: 314,
        posterUrl: null,
        tmdbId: 2,
        genres: [
          { code: "action", name: "アクション" },
          { code: "sf", name: "SF" },
        ],
      },
    ]);
  });
});
