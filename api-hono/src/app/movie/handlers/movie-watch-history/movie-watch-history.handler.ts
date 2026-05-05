import { AppType } from "@/core/app";
import { MovieIsNotFound } from "@/core/errors";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import {
  getMovieWatchHistory,
  WatchHistoryItem,
} from "../../services/watch-history/watch-history.service";

const paramSchema = z.object({
  movieId: z.coerce.number().int().positive(),
});

export default function (app: AppType) {
  app.get(
    "/movies/:movieId/watch-history",

    validator("param", (value, c) => {
      const parsed = paramSchema.safeParse(value);
      if (!parsed.success) {
        return c.json({ error: parsed.error }, StatusCodes.BAD_REQUEST);
      }
      return parsed.data;
    }),

    requireAuthMiddleware,

    async (c) => {
      const logger = c.get("logger");
      const operator = c.get("operator")!;
      const { movieId } = c.req.valid("param");

      const result = await getMovieWatchHistory(logger, operator.id, movieId);

      if (result.isErr()) {
        if (result.error instanceof MovieIsNotFound) {
          return c.json({ error: "movie not found" }, StatusCodes.NOT_FOUND);
        }
      }

      return c.json(
        result._unsafeUnwrap() satisfies WatchHistoryItem[],
        StatusCodes.OK,
      );
    },
  );
}
