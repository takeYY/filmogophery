import { environment } from "@/core/environment";
import { drizzle } from "drizzle-orm/mysql2";
import { createPool } from "mysql2/promise";

export const testDbPool = createPool({
  host: environment.WRITER_DATABASE.HOST,
  port: environment.WRITER_DATABASE.PORT,
  user: environment.WRITER_DATABASE.USER,
  password: environment.WRITER_DATABASE.PASSWORD,
  database: environment.WRITER_DATABASE.NAME,
});

export const testDb = drizzle(testDbPool);

export async function cleanupTables(tables: string[]) {
  const connection = await testDbPool.getConnection();
  try {
    await connection.query("SET FOREIGN_KEY_CHECKS = 0");
    for (const table of tables) {
      await connection.query(`TRUNCATE TABLE ${table}`);
    }
    await connection.query("SET FOREIGN_KEY_CHECKS = 1");
  } finally {
    connection.release();
  }
}
