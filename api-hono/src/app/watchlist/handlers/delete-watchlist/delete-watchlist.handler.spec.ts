import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import pino from "pino";
import * as deleteWatchlistService from "../../services/delete-watchlist/delete-watchlist.service";
import deleteWatchlistHandler from "./delete-watchlist.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

describe("delete-watchlist.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(deleteWatchlistService, "deleteWatchlist").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    deleteWatchlistHandler(testApp);
    return testApp;
  };

  test("should delete watchlist item and return 204", async () => {
    spyOn(deleteWatchlistService, "deleteWatchlist").mockResolvedValue(
      ok(undefined),
    );

    const res = await makeApp().request("/v1/watchlist/1", {
      method: "DELETE",
    });

    expect(res.status).toBe(StatusCodes.NO_CONTENT);
  });

  test("should return 404 when watchlist item not found", async () => {
    spyOn(deleteWatchlistService, "deleteWatchlist").mockResolvedValue(
      err(new deleteWatchlistService.WatchlistNotFoundError()),
    );

    const res = await makeApp().request("/v1/watchlist/999", {
      method: "DELETE",
    });

    expect(res.status).toBe(StatusCodes.NOT_FOUND);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/watchlist/1", {
      method: "DELETE",
    });

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
