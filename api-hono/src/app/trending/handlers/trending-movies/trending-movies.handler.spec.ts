import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import pino from "pino";
import * as trendingService from "../../services/trending-movies/trending-movies.service";
import trendingHandler from "./trending-movies.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

const mockTrendingMovies = [
  {
    id: 1,
    title: "Trending Movie 1",
    posterUrl: "/poster1.jpg",
    tmdbId: 101,
    hasReview: false,
  },
  {
    id: 2,
    title: "Trending Movie 2",
    posterUrl: "/poster2.jpg",
    tmdbId: 102,
    hasReview: true,
  },
];

describe("trending-movies.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(trendingService, "getTrendingMovies").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    trendingHandler(testApp);
    return testApp;
  };

  test("should return trending movies with hasReview flags", async () => {
    spyOn(trendingService, "getTrendingMovies").mockResolvedValue(
      ok(mockTrendingMovies),
    );

    const res = await makeApp().request("/v1/trending/movies");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockTrendingMovies);
  });

  test("should return empty array when no trending movies", async () => {
    spyOn(trendingService, "getTrendingMovies").mockResolvedValue(ok([]));

    const res = await makeApp().request("/v1/trending/movies");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual([]);
  });

  test("should return 500 when service fails", async () => {
    spyOn(trendingService, "getTrendingMovies").mockResolvedValue(
      err(new trendingService.TrendingMoviesError("tmdb error")),
    );

    const res = await makeApp().request("/v1/trending/movies");

    expect(res.status).toBe(StatusCodes.INTERNAL_SERVER_ERROR);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/trending/movies");

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
