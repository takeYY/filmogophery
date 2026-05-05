import { dbConnections } from "@/core/db";
import { genres } from "@/core/drizzle/schema";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function fetchAllGenres(
  db: MySql2Database = dbConnections.readonly,
) {
  return db.select({ code: genres.code, name: genres.name }).from(genres);
}
