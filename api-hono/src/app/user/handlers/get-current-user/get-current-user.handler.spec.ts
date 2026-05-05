import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import getCurrentUserHandler from "./get-current-user.handler";

const mockUser = {
  id: 1,
  username: "testuser",
  email: "test@example.com",
  passwordHash: "hash",
  isActive: 1,
  createdAt: "2024-01-01T00:00:00Z",
  updatedAt: "2024-01-01T00:00:00Z",
};

describe("get-current-user.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    getCurrentUserHandler(testApp);
    return testApp;
  };

  test("should return current user id and username", async () => {
    const res = await makeApp().request("/v1/users/me");

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual({ id: 1, username: "testuser" });
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/users/me");

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
