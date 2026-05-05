import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import * as platformsRepository from "../../repositories/platforms.repository";
import getPlatformsHandler from "./get-platforms.handler";

const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

const mockPlatforms = [
  { id: 1, code: "netflix", name: "Netflix" },
  { id: 2, code: "primeVideo", name: "Prime Video" },
];

describe("get-platforms.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(platformsRepository, "fetchAllPlatforms").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    getPlatformsHandler(testApp);
    return testApp;
  };

  test("should return platforms list", async () => {
    spyOn(platformsRepository, "fetchAllPlatforms").mockResolvedValue(
      mockPlatforms,
    );

    const res = await makeApp().request("/v1/platforms");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockPlatforms);
  });

  test("should return empty array when no platforms", async () => {
    spyOn(platformsRepository, "fetchAllPlatforms").mockResolvedValue([]);

    const res = await makeApp().request("/v1/platforms");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual([]);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/platforms");

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
