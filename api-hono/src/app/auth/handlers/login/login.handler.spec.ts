import { Variables } from "@/core/app";
import { afterEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import pino from "pino";
import * as loginService from "../../services/login/login.service";
import loginHandler from "./login.handler";

const mockToken = {
  accessToken: "access.token.value",
  refreshToken: "refreshtokenvalue",
  tokenType: "Bearer",
  expiresIn: 3600,
  expiresAt: "2026-04-28T13:00:00.000Z",
};

const testLogger = pino({ level: "silent" });

function createTestApp() {
  const app = new Hono<{ Variables: Variables }>().basePath("/v1");
  app.use(async (c, next) => {
    c.set("logger", testLogger);
    await next();
  });
  loginHandler(app);
  return app;
}

async function postLogin(app: ReturnType<typeof createTestApp>, body: unknown) {
  return app.request("/v1/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
}

describe("login.handler", () => {
  afterEach(() => {
    spyOn(loginService, "login").mockRestore();
  });

  test("正しい認証情報でトークンを返す", async () => {
    spyOn(loginService, "login").mockResolvedValue(ok(mockToken));

    const res = await postLogin(createTestApp(), {
      email: "test@example.com",
      password: "Password1",
    });

    expect(res.status).toBe(StatusCodes.OK);
    expect(await res.json()).toEqual(mockToken);
  });

  test("認証失敗時は 401 を返す", async () => {
    spyOn(loginService, "login").mockResolvedValue(
      err(
        new loginService.InvalidCredentialsError("invalid email or password"),
      ),
    );

    const res = await postLogin(createTestApp(), {
      email: "test@example.com",
      password: "Password1",
    });

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
    expect((await res.json()).message).toBe("invalid email or password");
  });

  test("システムエラー時は 500 を返す", async () => {
    spyOn(loginService, "login").mockResolvedValue(
      err(new Error("unexpected error")),
    );

    const res = await postLogin(createTestApp(), {
      email: "test@example.com",
      password: "Password1",
    });

    expect(res.status).toBe(StatusCodes.INTERNAL_SERVER_ERROR);
    expect((await res.json()).message).toBe("system error");
  });

  test("email が不正な場合は 400 を返す", async () => {
    const res = await postLogin(createTestApp(), {
      email: "not-an-email",
      password: "Password1",
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("password が短すぎる場合は 400 を返す", async () => {
    const res = await postLogin(createTestApp(), {
      email: "test@example.com",
      password: "short",
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("password に大文字・小文字・数字が含まれない場合は 400 を返す", async () => {
    const res = await postLogin(createTestApp(), {
      email: "test@example.com",
      password: "alllowercase",
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("リクエストボディが空の場合は 400 を返す", async () => {
    const res = await postLogin(createTestApp(), {});

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });
});
