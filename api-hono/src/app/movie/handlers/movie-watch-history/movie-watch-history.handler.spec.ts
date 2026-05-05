import { Variables } from "@/core/app";
import { MovieIsNotFound } from "@/core/errors";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import pino from "pino";
import * as watchHistoryService from "../../services/watch-history/watch-history.service";
import movieWatchHistoryHandler from "./movie-watch-history.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

const mockWatchHistory = [
  {
    id: 1,
    platform: { id: 1, code: "netflix", name: "Netflix" },
    watchedAt: "2025-01-01",
  },
  {
    id: 2,
    platform: { id: 2, code: "primeVideo", name: "Prime Video" },
    watchedAt: null,
  },
];

describe("movie-watch-history.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(watchHistoryService, "getMovieWatchHistory").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    movieWatchHistoryHandler(testApp);
    return testApp;
  };

  test("should return watch history list", async () => {
    spyOn(watchHistoryService, "getMovieWatchHistory").mockResolvedValue(
      ok(mockWatchHistory),
    );

    const res = await makeApp().request("/v1/movies/1/watch-history");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockWatchHistory);
  });

  test("should return empty array when no watch history", async () => {
    spyOn(watchHistoryService, "getMovieWatchHistory").mockResolvedValue(
      ok([]),
    );

    const res = await makeApp().request("/v1/movies/1/watch-history");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual([]);
  });

  test("should return 404 when movie not found", async () => {
    spyOn(watchHistoryService, "getMovieWatchHistory").mockResolvedValue(
      err(new MovieIsNotFound()),
    );

    const res = await makeApp().request("/v1/movies/999/watch-history");

    expect(res.status).toBe(StatusCodes.NOT_FOUND);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/movies/1/watch-history");

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
