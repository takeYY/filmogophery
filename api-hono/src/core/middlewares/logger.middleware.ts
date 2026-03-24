import { createMiddleware } from "hono/factory";
import { pinoLogger, Variables } from "../app";

export const loggerMiddleware = createMiddleware<{ Variables: Variables }>(
  async (c, next) => {
    const requestLogger = pinoLogger.child({ requestId: crypto.randomUUID() });
    c.set("logger", requestLogger);
    const start = Date.now();
    await next();
    requestLogger.info(
      { durationMs: Date.now() - start, status: c.res.status },
      "request completed",
    );
  },
);
