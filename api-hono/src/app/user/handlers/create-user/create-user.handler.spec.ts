import * as createUserService from "@/app/user/services/create-user/create-user.service";
import { UserAlreadyExistsError } from "@/app/user/services/create-user/create-user.service";
import { Variables } from "@/core/app";
import { afterEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import createUserHandler from "./create-user.handler";

const mockToken = {
  accessToken: "access.token.value",
  refreshToken: "refreshtokenvalue",
  tokenType: "Bearer",
  expiresIn: 3600,
  expiresAt: "2024-01-01T01:00:00.000Z",
};

const validBody = {
  username: "testuser",
  email: "test@example.com",
  password: "Password1",
};

function createTestApp() {
  const app = new Hono<{ Variables: Variables }>().basePath("/v1");
  createUserHandler(app);
  return app;
}

describe("create-user.handler", () => {
  afterEach(() => {
    spyOn(createUserService, "createUser").mockRestore();
  });

  test("正常にユーザーを作成して201を返す", async () => {
    spyOn(createUserService, "createUser").mockResolvedValue(ok(mockToken));

    const res = await createTestApp().request("/v1/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(validBody),
    });

    expect(res.status).toBe(StatusCodes.CREATED);
    expect(await res.json()).toEqual(mockToken);
  });

  test("username または email が重複している場合は409を返す", async () => {
    spyOn(createUserService, "createUser").mockResolvedValue(
      err(new UserAlreadyExistsError("username or email is already taken")),
    );

    const res = await createTestApp().request("/v1/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(validBody),
    });

    expect(res.status).toBe(StatusCodes.CONFLICT);
    const body = await res.json();
    expect(body.message).toBe("conflict");
  });

  test("予期しないエラーが発生した場合は500を返す", async () => {
    spyOn(createUserService, "createUser").mockResolvedValue(
      err(new Error("connection refused")),
    );

    const res = await createTestApp().request("/v1/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(validBody),
    });

    expect(res.status).toBe(StatusCodes.INTERNAL_SERVER_ERROR);
  });

  test("username が空の場合は400を返す", async () => {
    const res = await createTestApp().request("/v1/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...validBody, username: "" }),
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("email が不正な形式の場合は400を返す", async () => {
    const res = await createTestApp().request("/v1/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...validBody, email: "not-an-email" }),
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("password が8文字未満の場合は400を返す", async () => {
    const res = await createTestApp().request("/v1/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...validBody, password: "Pass1" }),
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("password に大文字が含まれない場合は400を返す", async () => {
    const res = await createTestApp().request("/v1/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...validBody, password: "password1" }),
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("password に数字が含まれない場合は400を返す", async () => {
    const res = await createTestApp().request("/v1/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ...validBody, password: "Passwordonly" }),
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });
});
