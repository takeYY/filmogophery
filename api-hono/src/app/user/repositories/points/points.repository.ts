import { dbConnections } from "@/core/db";
import { pointHistory, userPoints } from "@/core/drizzle/schema";
import { desc, eq } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function findUserPoints(
  userId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const [row] = await db
    .select()
    .from(userPoints)
    .where(eq(userPoints.userId, userId))
    .limit(1);
  return row ?? null;
}

export async function findPointHistory(
  userId: number,
  limit: number,
  offset: number,
  db: MySql2Database = dbConnections.readonly,
) {
  return db
    .select()
    .from(pointHistory)
    .where(eq(pointHistory.userId, userId))
    .orderBy(desc(pointHistory.createdAt))
    .limit(limit)
    .offset(offset);
}
