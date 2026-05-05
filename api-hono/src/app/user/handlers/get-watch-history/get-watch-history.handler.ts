import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import {
  getWatchHistory,
  WatchHistoryItem,
} from "../../services/get-watch-history/get-watch-history.service";

const querySchema = z.object({
  limit: z.coerce.number().int().positive().default(12),
  offset: z.coerce.number().int().min(0).default(0),
});

export default function (app: AppType) {
  app.get(
    "/users/me/watch-history",

    validator("query", (value, c) => {
      const parsed = querySchema.safeParse(value);
      if (!parsed.success) {
        return c.json({ error: parsed.error }, StatusCodes.BAD_REQUEST);
      }
      return parsed.data;
    }),

    requireAuthMiddleware,

    async (c) => {
      const logger = c.get("logger");
      const operator = c.get("operator")!;
      const query = c.req.valid("query");

      const result = await getWatchHistory(
        logger,
        operator.id,
        query.limit,
        query.offset,
      );

      return c.json(result.value satisfies WatchHistoryItem[], StatusCodes.OK);
    },
  );
}
