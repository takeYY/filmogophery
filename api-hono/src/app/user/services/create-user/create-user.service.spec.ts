import {
  createUser,
  UserAlreadyExistsError,
} from "@/app/user/services/create-user/create-user.service";
import * as dbModule from "@/core/db";
import { afterEach, describe, expect, mock, spyOn, test } from "bun:test";
import pino from "pino";

const logger = pino({ level: "silent" });

describe("createUser", () => {
  afterEach(() => {
    mock.restore();
  });

  test("正常にユーザーを作成してトークンを返す", async () => {
    mock.module("@/app/user/repositories/users/users.repository", () => ({
      insertUser: () => Promise.resolve(1),
    }));
    spyOn(dbModule.dbConnections.default, "insert").mockReturnValue({
      values: () => Promise.resolve(),
    } as any);

    const result = await createUser(
      logger,
      "testuser",
      "test@example.com",
      "Password1",
    );

    expect(result.isOk()).toBe(true);
    const token = result._unsafeUnwrap();
    expect(token.tokenType).toBe("Bearer");
    expect(token.expiresIn).toBe(3600);
    expect(token.accessToken).toBeString();
    expect(token.refreshToken).toBeString();
    expect(token.expiresAt).toBeString();
  });

  test("username または email が重複している場合は UserAlreadyExistsError を返す", async () => {
    mock.module("@/app/user/repositories/users/users.repository", () => ({
      insertUser: () =>
        Promise.reject(
          new Error("Duplicate entry 'testuser' for key 'username'"),
        ),
    }));

    const result = await createUser(
      logger,
      "testuser",
      "test@example.com",
      "Password1",
    );

    expect(result.isErr()).toBe(true);
    expect(result._unsafeUnwrapErr()).toBeInstanceOf(UserAlreadyExistsError);
  });

  test("予期しないエラーが発生した場合は Error を返す", async () => {
    mock.module("@/app/user/repositories/users/users.repository", () => ({
      insertUser: () => Promise.reject(new Error("connection refused")),
    }));

    const result = await createUser(
      logger,
      "testuser",
      "test@example.com",
      "Password1",
    );

    expect(result.isErr()).toBe(true);
    expect(result._unsafeUnwrapErr()).toBeInstanceOf(Error);
    expect(result._unsafeUnwrapErr()).not.toBeInstanceOf(
      UserAlreadyExistsError,
    );
  });
});
