const server = {
  PORT: process.env.SERVER_PORT,
};

const logger = {
  LEVEL: process.env.LOG_LEVEL,
  FORMAT: process.env.LOG_FORMAT,
};

const readerDatabase = {
  DB_CORE: process.env.READER_DB_CORE_COUNT,
  HOST: process.env.READER_DB_HOST?.split(":")[0],
  PORT: Number(process.env.READER_DB_PORT) ?? 3306,
  USER: process.env.READER_DB_USER,
  PASSWORD: process.env.READER_DB_PWD,
  NAME: process.env.READER_DB_NAME,
};

const writerDatabase = {
  DB_CORE: process.env.WRITER_DB_CORE_COUNT,
  HOST: process.env.WRITER_DB_HOST?.split(":")[0],
  PORT: Number(process.env.WRITER_DB_PORT) ?? 3306,
  USER: process.env.WRITER_DB_USER,
  PASSWORD: process.env.WRITER_DB_PWD,
  NAME: process.env.WRITER_DB_NAME,
};

const redis = {
  HOST: process.env.REDIS_HOST,
  PORT: process.env.REDIS_PORT,
  PASSWORD: process.env.REDIS_PASSWORD,
  DB: process.env.REDIS_DB,
};

const tmdb = {
  ACCESS_TOKEN: process.env.TMDB_ACCESS_TOKEN,
};

const token = {
  JWT_SECRET: process.env.JWT_SECRET,
};

export const environment = {
  SERVER: server,
  LOGGER: logger,
  READER_DATABASE: readerDatabase,
  WRITER_DATABASE: writerDatabase,
  REDIS: redis,
  TMDB: tmdb,
  TOKEN: token,
} as const;
