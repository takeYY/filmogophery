import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import * as genresRepository from "../../repositories/genres.repository";
import getGenresHandler from "./get-genres.handler";

const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

const mockGenres = [
  { code: "action", name: "アクション" },
  { code: "sf", name: "SF" },
];

describe("get-genres.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(genresRepository, "fetchAllGenres").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    getGenresHandler(testApp);
    return testApp;
  };

  test("should return genres list", async () => {
    spyOn(genresRepository, "fetchAllGenres").mockResolvedValue(mockGenres);

    const res = await makeApp().request("/v1/genres");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockGenres);
  });

  test("should return empty array when no genres", async () => {
    spyOn(genresRepository, "fetchAllGenres").mockResolvedValue([]);

    const res = await makeApp().request("/v1/genres");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual([]);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/genres");

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
