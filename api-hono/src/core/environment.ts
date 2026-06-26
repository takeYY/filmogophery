import { z } from "zod";

const envSchema = z.object({
  SERVER_PORT: z.string().min(1),
  LOG_LEVEL: z.string().default("info"),
  LOG_FORMAT: z.string().optional(),
  READER_DB_CORE_COUNT: z.coerce.number().optional(),
  READER_DB_HOST: z.string().min(1),
  READER_DB_PORT: z.coerce.number().default(3306),
  READER_DB_USER: z.string().min(1),
  READER_DB_PWD: z.string().min(1),
  READER_DB_NAME: z.string().min(1),
  WRITER_DB_CORE_COUNT: z.coerce.number().optional(),
  WRITER_DB_HOST: z.string().min(1),
  WRITER_DB_PORT: z.coerce.number().default(3306),
  WRITER_DB_USER: z.string().min(1),
  WRITER_DB_PWD: z.string().min(1),
  WRITER_DB_NAME: z.string().min(1),
  REDIS_HOST: z.string().min(1),
  REDIS_PORT: z.string().min(1),
  REDIS_PASSWORD: z.string().optional(),
  REDIS_DB: z.string().optional(),
  TMDB_ACCESS_TOKEN: z.string().min(1),
  JWT_SECRET: z.string().min(1),
});

const parsed = envSchema.safeParse(process.env);

if (!parsed.success) {
  console.error("環境変数が不正です:", parsed.error.flatten().fieldErrors);
  process.exit(1);
}

const env = parsed.data;

export const environment = {
  SERVER: {
    PORT: env.SERVER_PORT,
  },
  LOGGER: {
    LEVEL: env.LOG_LEVEL,
    FORMAT: env.LOG_FORMAT,
  },
  READER_DATABASE: {
    DB_CORE: env.READER_DB_CORE_COUNT,
    HOST: env.READER_DB_HOST.split(":")[0],
    PORT: env.READER_DB_PORT,
    USER: env.READER_DB_USER,
    PASSWORD: env.READER_DB_PWD,
    NAME: env.READER_DB_NAME,
  },
  WRITER_DATABASE: {
    DB_CORE: env.WRITER_DB_CORE_COUNT,
    HOST: env.WRITER_DB_HOST.split(":")[0],
    PORT: env.WRITER_DB_PORT,
    USER: env.WRITER_DB_USER,
    PASSWORD: env.WRITER_DB_PWD,
    NAME: env.WRITER_DB_NAME,
  },
  REDIS: {
    HOST: env.REDIS_HOST,
    PORT: env.REDIS_PORT,
    PASSWORD: env.REDIS_PASSWORD,
    DB: env.REDIS_DB,
  },
  TMDB: {
    ACCESS_TOKEN: env.TMDB_ACCESS_TOKEN,
  },
  TOKEN: {
    JWT_SECRET: env.JWT_SECRET,
  },
} as const;
