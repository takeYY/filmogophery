import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import pino from "pino";
import * as createWatchHistoryService from "../../services/create-watch-history/create-watch-history.service";
import createWatchHistoryHandler from "./create-watch-history.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

describe("create-watch-history.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(createWatchHistoryService, "createWatchHistory").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    createWatchHistoryHandler(testApp);
    return testApp;
  };

  const postHistory = (
    app: ReturnType<typeof makeApp>,
    id: number,
    body: unknown,
  ) =>
    app.request(`/v1/reviews/${id}/history`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

  test("should create watch history with platformId only", async () => {
    spyOn(createWatchHistoryService, "createWatchHistory").mockResolvedValue(
      ok(undefined),
    );

    const res = await postHistory(makeApp(), 1, { platformId: 1 });

    expect(res.status).toBe(StatusCodes.CREATED);
  });

  test("should create watch history with platformId and watchedDate", async () => {
    spyOn(createWatchHistoryService, "createWatchHistory").mockResolvedValue(
      ok(undefined),
    );

    const res = await postHistory(makeApp(), 1, {
      platformId: 1,
      watchedDate: "2025-01-01",
    });

    expect(res.status).toBe(StatusCodes.CREATED);
  });

  test("should return 400 when platformId is missing", async () => {
    const res = await postHistory(makeApp(), 1, {});

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 400 when watchedDate is in the future", async () => {
    const res = await postHistory(makeApp(), 1, {
      platformId: 1,
      watchedDate: "2099-12-31",
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 404 when review not found", async () => {
    spyOn(createWatchHistoryService, "createWatchHistory").mockResolvedValue(
      err(new createWatchHistoryService.ReviewNotFoundError()),
    );

    const res = await postHistory(makeApp(), 999, { platformId: 1 });

    expect(res.status).toBe(StatusCodes.NOT_FOUND);
  });

  test("should return 404 when platform not found", async () => {
    spyOn(createWatchHistoryService, "createWatchHistory").mockResolvedValue(
      err(new createWatchHistoryService.PlatformNotFoundError()),
    );

    const res = await postHistory(makeApp(), 1, { platformId: 999 });

    expect(res.status).toBe(StatusCodes.NOT_FOUND);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await postHistory(makeApp(), 1, { platformId: 1 });

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
