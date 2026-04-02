import { Hono } from "hono";
import pino, { Logger } from "pino";
import { users } from "./drizzle/schema";
import { loggerMiddleware } from "./middlewares/logger.middleware";

export const pinoLogger = pino({ level: process.env.LOG_LEVEL ?? "info" });

export type Variables = {
  logger: Logger;
  operator?: typeof users.$inferSelect;
};

export const app = new Hono<{ Variables: Variables }>().basePath("/v1");

app.use(loggerMiddleware);

export type AppType = typeof app;
