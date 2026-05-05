import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { ok } from "neverthrow";
import pino from "pino";
import * as moviesService from "../../services/movies/movies.service";
import moviesHandler from "./movies.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

describe("movies.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(moviesService, "getMovies").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    moviesHandler(testApp);
    return testApp;
  };

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

    const res = await makeApp().request("/v1/movies");

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

    const res = await makeApp().request(
      "/v1/movies?genre=action&limit=10&offset=0",
    );

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockMovies);
  });

  test("should return empty array when no movies found", async () => {
    spyOn(moviesService, "getMovies").mockResolvedValue(ok([]));

    const res = await makeApp().request("/v1/movies");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual([]);
  });

  test("should return 400 for invalid limit", async () => {
    const res = await makeApp().request("/v1/movies?limit=-1");

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 400 for invalid offset", async () => {
    const res = await makeApp().request("/v1/movies?offset=-1");

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/movies");

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
