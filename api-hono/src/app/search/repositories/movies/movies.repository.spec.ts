import { genres, movieGenres, movies } from "@/core/drizzle/schema";
import { cleanupTables, testDb } from "@/core/test-helpers/db";
import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { batchCreateMovies, getMoviesByTmdbIds } from "./movies.repository";

describe("movies.repository", () => {
  beforeEach(async () => {
    await cleanupTables(["movie_genres", "movies", "genres"]);
  });

  afterEach(async () => {
    await cleanupTables(["movie_genres", "movies", "genres"]);
  });

  describe("getMoviesByTmdbIds", () => {
    test("should return empty array when no movies exist", async () => {
      const result = await getMoviesByTmdbIds([1, 2, 3], testDb);

      expect(result).toHaveLength(0);
    });

    test("should return empty array when no movies match the given tmdbIds", async () => {
      // テストデータ作成
      await testDb.insert(movies).values({
        tmdbId: 100,
        title: "Existing Movie",
        overview: "Overview",
        releaseDate: "2024-01-01",
        runtimeMinutes: 120,
      });

      // テスト実行（存在しないtmdbIdで検索）
      const result = await getMoviesByTmdbIds([999], testDb);

      // 検証
      expect(result).toHaveLength(0);
    });

    test("should return movies matching the given tmdbIds", async () => {
      // テストデータ作成
      const [genre1] = await testDb.insert(genres).values({
        code: "action",
        name: "Action",
      });

      const [movie1] = await testDb.insert(movies).values({
        tmdbId: 1,
        title: "Test Movie 1",
        overview: "Overview 1",
        releaseDate: "2024-01-01",
        runtimeMinutes: 120,
      });

      await testDb.insert(movieGenres).values({
        movieId: movie1.insertId,
        genreId: genre1.insertId,
      });

      // テスト実行
      const result = await getMoviesByTmdbIds([1], testDb);

      // 検証
      expect(result).toHaveLength(1);
      expect(result[0].tmdbId).toBe(1);
      expect(result[0].title).toBe("Test Movie 1");
    });

    test("should return multiple movies matching the given tmdbIds", async () => {
      // テストデータ作成
      const [genre1] = await testDb.insert(genres).values({
        code: "action",
        name: "Action",
      });

      const [movie1] = await testDb.insert(movies).values({
        tmdbId: 1,
        title: "Movie 1",
        overview: "Overview 1",
        releaseDate: "2024-01-01",
        runtimeMinutes: 120,
      });

      const [movie2] = await testDb.insert(movies).values({
        tmdbId: 2,
        title: "Movie 2",
        overview: "Overview 2",
        releaseDate: "2024-01-02",
        runtimeMinutes: 90,
      });

      await testDb.insert(movieGenres).values([
        { movieId: movie1.insertId, genreId: genre1.insertId },
        { movieId: movie2.insertId, genreId: genre1.insertId },
      ]);

      // テスト実行
      const result = await getMoviesByTmdbIds([1, 2], testDb);

      // 検証
      expect(result).toHaveLength(2);
      const tmdbIds = result.map((m) => m.tmdbId);
      expect(tmdbIds).toContain(1);
      expect(tmdbIds).toContain(2);
    });

    test("should include genre codes and names in the result", async () => {
      // テストデータ作成
      const [genre1] = await testDb.insert(genres).values({
        code: "action",
        name: "Action",
      });

      const [movie1] = await testDb.insert(movies).values({
        tmdbId: 1,
        title: "Action Movie",
        overview: "Overview",
        releaseDate: "2024-01-01",
        runtimeMinutes: 120,
      });

      await testDb.insert(movieGenres).values({
        movieId: movie1.insertId,
        genreId: genre1.insertId,
      });

      // テスト実行
      const result = await getMoviesByTmdbIds([1], testDb);

      // 検証
      expect(result).toHaveLength(1);
      expect(result[0].genreCodes).toBe("action");
      expect(result[0].genreNames).toBe("Action");
    });
  });

  describe("batchCreateMovies", () => {
    test("should create movies and return their inserted ids", async () => {
      // テストデータ作成
      const [genre1] = await testDb.insert(genres).values({
        code: "action",
        name: "Action",
      });

      const newMovies = [
        {
          tmdbId: 10,
          title: "New Movie 1",
          overview: "Overview 1",
          releaseDate: "2024-01-01",
          runtimeMinutes: 120,
          posterUrl: null,
          genreIds: [genre1.insertId],
        },
      ];

      // テスト実行
      const result = await batchCreateMovies(newMovies, testDb);

      // 検証
      expect(result).toHaveLength(1);

      const created = await getMoviesByTmdbIds([10], testDb);
      expect(created).toHaveLength(1);
      expect(created[0].title).toBe("New Movie 1");
    });

    test("should create multiple movies in a single transaction", async () => {
      // テストデータ作成
      const [genre1] = await testDb.insert(genres).values({
        code: "drama",
        name: "Drama",
      });

      const newMovies = [
        {
          tmdbId: 20,
          title: "Batch Movie 1",
          overview: "Overview 1",
          releaseDate: "2024-01-01",
          runtimeMinutes: 100,
          posterUrl: null,
          genreIds: [genre1.insertId],
        },
        {
          tmdbId: 21,
          title: "Batch Movie 2",
          overview: "Overview 2",
          releaseDate: "2024-01-02",
          runtimeMinutes: 110,
          posterUrl: "https://example.com/poster.jpg",
          genreIds: [genre1.insertId],
        },
      ];

      // テスト実行
      const result = await batchCreateMovies(newMovies, testDb);

      // 検証
      expect(result).toHaveLength(2);

      const created = await getMoviesByTmdbIds([20, 21], testDb);
      expect(created).toHaveLength(2);
    });

    test("should create movie without genres when genreIds is empty", async () => {
      const newMovies = [
        {
          tmdbId: 30,
          title: "No Genre Movie",
          overview: "Overview",
          releaseDate: "2024-01-01",
          runtimeMinutes: 90,
          posterUrl: null,
          genreIds: [],
        },
      ];

      // テスト実行
      const result = await batchCreateMovies(newMovies, testDb);

      // 検証
      expect(result).toHaveLength(1);

      const created = await getMoviesByTmdbIds([30], testDb);
      expect(created).toHaveLength(1);
      expect(created[0].genreCodes).toBeNull();
    });
  });
});
