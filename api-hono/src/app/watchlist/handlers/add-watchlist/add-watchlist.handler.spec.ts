import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import pino from "pino";
import * as addWatchlistService from "../../services/add-watchlist/add-watchlist.service";
import addWatchlistHandler from "./add-watchlist.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

describe("add-watchlist.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(addWatchlistService, "addWatchlist").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    addWatchlistHandler(testApp);
    return testApp;
  };

  const post = (app: ReturnType<typeof makeApp>, body: unknown) =>
    app.request("/v1/watchlist", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

  test("should add to watchlist with movieId only (default priority)", async () => {
    spyOn(addWatchlistService, "addWatchlist").mockResolvedValue(ok(undefined));

    const res = await post(makeApp(), { movieId: 1 });

    expect(res.status).toBe(StatusCodes.CREATED);
  });

  test("should add to watchlist with movieId and priority", async () => {
    spyOn(addWatchlistService, "addWatchlist").mockResolvedValue(ok(undefined));

    const res = await post(makeApp(), { movieId: 1, priority: 3 });

    expect(res.status).toBe(StatusCodes.CREATED);
  });

  test("should return 400 when movieId is missing", async () => {
    const res = await post(makeApp(), {});

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 400 when priority is out of range", async () => {
    const res = await post(makeApp(), { movieId: 1, priority: 6 });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 404 when movie not found", async () => {
    spyOn(addWatchlistService, "addWatchlist").mockResolvedValue(
      err(new addWatchlistService.MovieNotFoundError()),
    );

    const res = await post(makeApp(), { movieId: 999 });

    expect(res.status).toBe(StatusCodes.NOT_FOUND);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await post(makeApp(), { movieId: 1 });

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
