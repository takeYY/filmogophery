import { users } from "@/core/drizzle/schema";
import { cleanupTables, testDb } from "@/core/test-helpers/db";
import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { findUserByEmail, updateLastLoginAt } from "./users.repository";

describe("users.repository", () => {
  beforeEach(async () => {
    await cleanupTables(["refresh_tokens", "users"]);
  });

  afterEach(async () => {
    await cleanupTables(["refresh_tokens", "users"]);
  });

  describe("findUserByEmail", () => {
    test("should return null when no user exists", async () => {
      const result = await findUserByEmail("notfound@example.com", testDb);

      expect(result).toBeNull();
    });

    test("should return the user matching the given email", async () => {
      // テストデータ作成
      await testDb.insert(users).values({
        username: "testuser",
        email: "test@example.com",
        passwordHash: "hashedpassword",
      });

      // テスト実行
      const result = await findUserByEmail("test@example.com", testDb);

      // 検証
      expect(result).not.toBeNull();
      expect(result!.email).toBe("test@example.com");
      expect(result!.username).toBe("testuser");
    });

    test("should return null when email does not match", async () => {
      // テストデータ作成
      await testDb.insert(users).values({
        username: "testuser",
        email: "test@example.com",
        passwordHash: "hashedpassword",
      });

      // テスト実行
      const result = await findUserByEmail("other@example.com", testDb);

      // 検証
      expect(result).toBeNull();
    });

    test("should return only the first user when multiple users exist", async () => {
      // テストデータ作成
      await testDb.insert(users).values([
        {
          username: "user1",
          email: "user1@example.com",
          passwordHash: "hash1",
        },
        {
          username: "user2",
          email: "user2@example.com",
          passwordHash: "hash2",
        },
      ]);

      // テスト実行
      const result = await findUserByEmail("user1@example.com", testDb);

      // 検証
      expect(result).not.toBeNull();
      expect(result!.email).toBe("user1@example.com");
    });
  });

  describe("updateLastLoginAt", () => {
    test("should update lastLoginAt for the given user", async () => {
      // テストデータ作成
      const [inserted] = await testDb.insert(users).values({
        username: "testuser",
        email: "test@example.com",
        passwordHash: "hashedpassword",
      });
      const userId = inserted.insertId as number;

      const lastLoginAt = "2026-04-27 12:00:00";

      // テスト実行
      await updateLastLoginAt(userId, lastLoginAt, testDb);

      // 検証
      const result = await findUserByEmail("test@example.com", testDb);
      expect(result!.lastLoginAt).toBe(lastLoginAt);
    });

    test("should not affect other users when updating lastLoginAt", async () => {
      // テストデータ作成
      const [inserted1] = await testDb.insert(users).values({
        username: "user1",
        email: "user1@example.com",
        passwordHash: "hash1",
      });
      await testDb.insert(users).values({
        username: "user2",
        email: "user2@example.com",
        passwordHash: "hash2",
      });

      const userId1 = inserted1.insertId as number;
      const lastLoginAt = "2026-04-27 12:00:00";

      // テスト実行
      await updateLastLoginAt(userId1, lastLoginAt, testDb);

      // 検証
      const user2 = await findUserByEmail("user2@example.com", testDb);
      expect(user2!.lastLoginAt).toBeNull();
    });
  });
});
