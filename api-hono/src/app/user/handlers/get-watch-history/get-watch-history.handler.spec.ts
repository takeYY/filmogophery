import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { ok } from "neverthrow";
import pino from "pino";
import * as watchHistoryService from "../../services/get-watch-history/get-watch-history.service";
import getWatchHistoryHandler from "./get-watch-history.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

const mockWatchHistory = [
  {
    id: 1,
    watchedAt: "2025-01-01",
    platform: { id: 1, code: "netflix", name: "Netflix" },
    movie: {
      id: 10,
      title: "テスト映画",
      overview: "概要",
      releaseDate: "2024-01-01",
      runtimeMinutes: 120,
      posterUrl: "/poster.jpg",
      tmdbId: 999,
      genres: [{ code: "action", name: "アクション" }],
    },
  },
];

describe("get-watch-history.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(watchHistoryService, "getWatchHistory").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    getWatchHistoryHandler(testApp);
    return testApp;
  };

  test("should return watch history list with default params", async () => {
    spyOn(watchHistoryService, "getWatchHistory").mockResolvedValue(
      ok(mockWatchHistory),
    );

    const res = await makeApp().request("/v1/users/me/watch-history");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockWatchHistory);
  });

  test("should return empty array when no watch history", async () => {
    spyOn(watchHistoryService, "getWatchHistory").mockResolvedValue(ok([]));

    const res = await makeApp().request("/v1/users/me/watch-history");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual([]);
  });

  test("should pass limit and offset to service", async () => {
    const spy = spyOn(watchHistoryService, "getWatchHistory").mockResolvedValue(
      ok([]),
    );

    await makeApp().request("/v1/users/me/watch-history?limit=5&offset=10");

    expect(spy).toHaveBeenCalledWith(expect.anything(), mockUser.id, 5, 10);
  });

  test("should return 400 for invalid limit", async () => {
    const res = await makeApp().request("/v1/users/me/watch-history?limit=-1");

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/users/me/watch-history");

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
