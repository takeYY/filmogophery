import { dbConnections } from "@/core/db";
import { platforms } from "@/core/drizzle/schema";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function fetchAllPlatforms(
  db: MySql2Database = dbConnections.readonly,
) {
  return db
    .select({ id: platforms.id, code: platforms.code, name: platforms.name })
    .from(platforms);
}
