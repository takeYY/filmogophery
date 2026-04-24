import { users } from "@/core/drizzle/schema";
import { cleanupTables, testDb } from "@/core/test-helpers/db";
import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { insertUser } from "./users.repository";

describe("users.repository", () => {
  beforeEach(async () => {
    await cleanupTables(["refresh_tokens", "users"]);
  });

  afterEach(async () => {
    await cleanupTables(["refresh_tokens", "users"]);
  });

  describe("insertUser", () => {
    test("should insert a user and return the inserted id", async () => {
      // テスト実行
      const insertedId = await insertUser(
        {
          username: "testuser",
          email: "test@example.com",
          passwordHash: "hashedpassword",
          lastLoginAt: "2024-01-01 00:00:00",
        },
        testDb,
      );

      // 検証
      expect(insertedId).toBeGreaterThan(0);

      const allUsers = await testDb.select().from(users);
      const user = allUsers.find((u) => u.id === insertedId);
      expect(user).toBeDefined();
      expect(user!.username).toBe("testuser");
      expect(user!.email).toBe("test@example.com");
    });

    test("should throw an error when inserting a duplicate username", async () => {
      // テストデータ作成
      await insertUser(
        {
          username: "duplicateuser",
          email: "first@example.com",
          passwordHash: "hashedpassword",
          lastLoginAt: "2024-01-01 00:00:00",
        },
        testDb,
      );

      // テスト実行・検証
      expect(
        insertUser(
          {
            username: "duplicateuser",
            email: "second@example.com",
            passwordHash: "hashedpassword",
            lastLoginAt: "2024-01-01 00:00:00",
          },
          testDb,
        ),
      ).rejects.toThrow();
    });

    test("should throw an error when inserting a duplicate email", async () => {
      // テストデータ作成
      await insertUser(
        {
          username: "user1",
          email: "duplicate@example.com",
          passwordHash: "hashedpassword",
          lastLoginAt: "2024-01-01 00:00:00",
        },
        testDb,
      );

      // テスト実行・検証
      expect(
        insertUser(
          {
            username: "user2",
            email: "duplicate@example.com",
            passwordHash: "hashedpassword",
            lastLoginAt: "2024-01-01 00:00:00",
          },
          testDb,
        ),
      ).rejects.toThrow();
    });
  });
});
