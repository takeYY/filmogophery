import { dbConnections } from "@/core/db";
import { refreshTokens, users } from "@/core/drizzle/schema";
import { and, eq, gt, isNull } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function findUserByEmail(
  email: string,
  db: MySql2Database = dbConnections.readonly,
) {
  const [user] = await db
    .select()
    .from(users)
    .where(eq(users.email, email))
    .limit(1);
  return user ?? null;
}

export async function updateLastLoginAt(
  userId: number,
  lastLoginAt: string,
  db: MySql2Database = dbConnections.default,
) {
  await db.update(users).set({ lastLoginAt }).where(eq(users.id, userId));
}

/**
 * ユーザーの有効なリフレッシュトークンをすべて無効化する
 */
export async function revokeActiveTokensByUserId(
  userId: number,
  now: string,
  db: MySql2Database = dbConnections.default,
) {
  await db
    .update(refreshTokens)
    .set({ revokedAt: now })
    .where(
      and(
        eq(refreshTokens.userId, userId),
        isNull(refreshTokens.revokedAt),
        gt(refreshTokens.expiresAt, now),
      ),
    );
}
