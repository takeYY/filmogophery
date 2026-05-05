import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import pino from "pino";
import * as usersRepository from "../../repositories/users/users.repository";
import logoutHandler from "./logout.handler";

const testLogger = pino({ level: "silent" });

const mockUser = {
  id: 1,
  username: "testuser",
  email: "test@example.com",
  passwordHash: "hash",
  isActive: 1,
  createdAt: "2024-01-01T00:00:00Z",
  updatedAt: "2024-01-01T00:00:00Z",
};

describe("logout.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
    spyOn(usersRepository, "revokeActiveTokensByUserId").mockResolvedValue(
      undefined,
    );
  });

  afterEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
    spyOn(usersRepository, "revokeActiveTokensByUserId").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    logoutHandler(testApp);
    return testApp;
  };

  test("should return 204 on successful logout", async () => {
    const res = await makeApp().request("/v1/auth/logout", {
      method: "POST",
    });

    expect(res.status).toBe(StatusCodes.NO_CONTENT);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/auth/logout", {
      method: "POST",
    });

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
