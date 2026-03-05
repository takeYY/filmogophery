import { afterEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { ok } from "neverthrow";
import * as moviesService from "../../services/movies/movies.service";
import moviesHandler from "./movies.handler";

describe("movies.handler", () => {
  afterEach(() => {
    spyOn(moviesService, "getMovies").mockRestore();
  });

  test("should return movies with default params", async () => {
    const mockMovies = [
      {
        id: 1,
        title: "Test Movie",
        overview: "Overview",
        releaseDate: "2024-01-01",
        runtimeMinute: 120,
        posterUrl: null,
        tmdbId: 1,
        genres: [{ code: "action", name: "Action" }],
      },
    ];

    spyOn(moviesService, "getMovies").mockResolvedValue(ok(mockMovies));

    const testApp = new Hono().basePath("/v1");
    moviesHandler(testApp);

    const res = await testApp.request("/v1/movies");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockMovies);
  });

  test("should return movies with genre filter", async () => {
    const mockMovies = [
      {
        id: 1,
        title: "Action Movie",
        overview: "Overview",
        releaseDate: "2024-01-01",
        runtimeMinute: 120,
        posterUrl: null,
        tmdbId: 1,
        genres: [{ code: "action", name: "Action" }],
      },
    ];

    spyOn(moviesService, "getMovies").mockResolvedValue(ok(mockMovies));

    const testApp = new Hono().basePath("/v1");
    moviesHandler(testApp);

    const res = await testApp.request(
      "/v1/movies?genre=action&limit=10&offset=0",
    );

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockMovies);
  });

  test("should return empty array when no movies found", async () => {
    spyOn(moviesService, "getMovies").mockResolvedValue(ok([]));

    const testApp = new Hono().basePath("/v1");
    moviesHandler(testApp);

    const res = await testApp.request("/v1/movies");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual([]);
  });

  test("should return 400 for invalid limit", async () => {
    const testApp = new Hono().basePath("/v1");
    moviesHandler(testApp);

    const res = await testApp.request("/v1/movies?limit=-1");

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 400 for invalid offset", async () => {
    const testApp = new Hono().basePath("/v1");
    moviesHandler(testApp);

    const res = await testApp.request("/v1/movies?offset=-1");

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });
});
