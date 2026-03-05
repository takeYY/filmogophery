import { drizzle } from "drizzle-orm/mysql2";
import { createPool } from "mysql2/promise";

import { environment } from "./environment";

const readerConnection = createPool({
  host: environment.READER_DATABASE.HOST,
  port: environment.READER_DATABASE.PORT,
  user: environment.READER_DATABASE.USER,
  password: environment.READER_DATABASE.PASSWORD,
  database: environment.READER_DATABASE.NAME,
});

const writerConnection = createPool({
  host: environment.WRITER_DATABASE.HOST,
  port: environment.WRITER_DATABASE.PORT,
  user: environment.WRITER_DATABASE.USER,
  password: environment.WRITER_DATABASE.PASSWORD,
  database: environment.WRITER_DATABASE.NAME,
});

export const dbConnections = {
  default: drizzle(writerConnection, { logger: Boolean(process.env.DB_LOG) }),
  readonly: drizzle(readerConnection, { logger: Boolean(process.env.DB_LOG) }),
} as const;

export type DbConnectionName = keyof typeof dbConnections;
