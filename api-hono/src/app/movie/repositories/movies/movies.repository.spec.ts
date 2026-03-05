import { genres, movieGenres, movies } from "@/core/drizzle/schema";
import { cleanupTables, testDb } from "@/core/test-helpers/db";
import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { getMoviesByGenre } from "./movies.repository";

describe("movies.repository", () => {
  beforeEach(async () => {
    await cleanupTables(["movie_genres", "movies", "genres"]);
  });

  afterEach(async () => {
    await cleanupTables(["movie_genres", "movies", "genres"]);
  });

  test("should return empty array when no movies exist", async () => {
    // テスト実行
    const result = await getMoviesByGenre(undefined, 10, 0, testDb);

    // 検証
    expect(result).toHaveLength(0);
  });

  test("should return empty array when no movies match the genre", async () => {
    // テストデータ作成
    const [actionGenre] = await testDb.insert(genres).values({
      code: "action",
      name: "Action",
    });

    const [movie1] = await testDb.insert(movies).values({
      tmdbId: 1,
      title: "Action Movie",
      overview: "Action overview",
      releaseDate: "2024-01-01",
      runtimeMinutes: 120,
    });

    await testDb.insert(movieGenres).values({
      movieId: movie1.insertId,
      genreId: actionGenre.insertId,
    });

    // テスト実行（存在しないジャンルで検索）
    const result = await getMoviesByGenre("drama", 10, 0, testDb);

    // 検証
    expect(result).toHaveLength(0);
  });

  test("should return movies without genre filter", async () => {
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
    const result = await getMoviesByGenre(undefined, 10, 0, testDb);

    // 検証
    expect(result).toHaveLength(1);
    expect(result[0].title).toBe("Test Movie 1");
  });

  test("should return movies filtered by genre", async () => {
    // テストデータ作成
    const [actionGenre] = await testDb.insert(genres).values({
      code: "action",
      name: "Action",
    });

    const [dramaGenre] = await testDb.insert(genres).values({
      code: "drama",
      name: "Drama",
    });

    const [movie1] = await testDb.insert(movies).values({
      tmdbId: 1,
      title: "Action Movie",
      overview: "Action overview",
      releaseDate: "2024-01-01",
      runtimeMinutes: 120,
    });

    const [movie2] = await testDb.insert(movies).values({
      tmdbId: 2,
      title: "Drama Movie",
      overview: "Drama overview",
      releaseDate: "2024-01-02",
      runtimeMinutes: 110,
    });

    await testDb.insert(movieGenres).values([
      { movieId: movie1.insertId, genreId: actionGenre.insertId },
      { movieId: movie2.insertId, genreId: dramaGenre.insertId },
    ]);

    // テスト実行
    const result = await getMoviesByGenre("action", 10, 0, testDb);

    // 検証
    expect(result).toHaveLength(1);
    expect(result[0].title).toBe("Action Movie");
  });
});
