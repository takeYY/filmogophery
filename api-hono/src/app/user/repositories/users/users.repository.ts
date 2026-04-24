import { dbConnections } from "@/core/db";
import { users } from "@/core/drizzle/schema";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function insertUser(
  user: {
    username: string;
    email: string;
    passwordHash: string;
    lastLoginAt?: string;
  },
  db: MySql2Database = dbConnections.default,
) {
  const [result] = await db.insert(users).values(user);
  return result.insertId as number;
}
