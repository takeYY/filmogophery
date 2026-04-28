import * as usersRepo from "@/app/auth/repositories/users/users.repository";
import { afterEach, describe, expect, mock, spyOn, test } from "bun:test";
import pino from "pino";
import { InvalidCredentialsError, login } from "./login.service";

const logger = pino({ level: "silent" });

const mockInsertRefreshToken = async () => {};

// ユーザーのモックデータ
const mockUser = {
  id: 1,
  username: "testuser",
  email: "test@example.com",
  passwordHash: await Bun.password.hash("Password1", { algorithm: "bcrypt" }),
  isActive: 1,
  lastLoginAt: null,
  createdAt: "2026-01-01 00:00:00",
  updatedAt: "2026-01-01 00:00:00",
};

describe("login", () => {
  afterEach(() => {
    mock.restore();
  });

  test("メールアドレスが存在しない場合は InvalidCredentialsError を返す", async () => {
    spyOn(usersRepo, "findUserByEmail").mockResolvedValue(null as any);

    const result = await login(logger, "notfound@example.com", "Password1");

    expect(result.isErr()).toBe(true);
    expect(result._unsafeUnwrapErr()).toBeInstanceOf(InvalidCredentialsError);
  });

  test("パスワードが一致しない場合は InvalidCredentialsError を返す", async () => {
    spyOn(usersRepo, "findUserByEmail").mockResolvedValue(mockUser);

    const result = await login(logger, "test@example.com", "WrongPassword1");

    expect(result.isErr()).toBe(true);
    expect(result._unsafeUnwrapErr()).toBeInstanceOf(InvalidCredentialsError);
  });

  test("認証成功時はトークン情報を返す", async () => {
    spyOn(usersRepo, "findUserByEmail").mockResolvedValue(mockUser);
    spyOn(usersRepo, "updateLastLoginAt").mockResolvedValue(undefined);

    const result = await login(logger, "test@example.com", "Password1", {
      insertRefreshToken: mockInsertRefreshToken,
    });

    expect(result.isOk()).toBe(true);
    const token = result._unsafeUnwrap();
    expect(token.tokenType).toBe("Bearer");
    expect(token.expiresIn).toBe(3600);
    expect(token.accessToken).toBeString();
    expect(token.refreshToken).toBeString();
    expect(token.expiresAt).toBeString();
  });

  test("認証成功時は lastLoginAt が更新される", async () => {
    spyOn(usersRepo, "findUserByEmail").mockResolvedValue(mockUser);
    const updateSpy = spyOn(usersRepo, "updateLastLoginAt").mockResolvedValue(
      undefined,
    );

    await login(logger, "test@example.com", "Password1", {
      insertRefreshToken: mockInsertRefreshToken,
    });

    expect(updateSpy).toHaveBeenCalledTimes(1);
    expect(updateSpy).toHaveBeenCalledWith(mockUser.id, expect.any(String));
  });

  test("認証成功時はリフレッシュトークンが保存される", async () => {
    spyOn(usersRepo, "findUserByEmail").mockResolvedValue(mockUser);
    spyOn(usersRepo, "updateLastLoginAt").mockResolvedValue(undefined);

    let capturedArgs: unknown[] = [];
    const insertSpy = async (...args: unknown[]) => {
      capturedArgs = args;
    };

    await login(logger, "test@example.com", "Password1", {
      insertRefreshToken: insertSpy as typeof mockInsertRefreshToken,
    });

    expect(capturedArgs[0]).toBe(mockUser.id);
    expect(typeof capturedArgs[1]).toBe("string"); // tokenHash
    expect(typeof capturedArgs[2]).toBe("string"); // expiresAt
  });

  test("DB 保存でエラーが発生した場合は err を返す", async () => {
    spyOn(usersRepo, "findUserByEmail").mockResolvedValue(mockUser);

    const result = await login(logger, "test@example.com", "Password1", {
      insertRefreshToken: async () => {
        throw new Error("DB error");
      },
    });

    expect(result.isErr()).toBe(true);
    expect(result._unsafeUnwrapErr()).toBeInstanceOf(Error);
  });
});
