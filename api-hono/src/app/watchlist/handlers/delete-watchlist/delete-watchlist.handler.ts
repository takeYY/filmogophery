import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import {
  deleteWatchlist,
  WatchlistNotFoundError,
} from "../../services/delete-watchlist/delete-watchlist.service";

const paramSchema = z.object({
  watchlistId: z.coerce.number().int().positive(),
});

export default function (app: AppType) {
  app.delete(
    "/watchlist/:watchlistId",

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
      const { watchlistId } = c.req.valid("param");

      const result = await deleteWatchlist(logger, watchlistId);

      if (result.isErr()) {
        if (result.error instanceof WatchlistNotFoundError) {
          return c.json(
            { error: "watchlist not found" },
            StatusCodes.NOT_FOUND,
          );
        }
      }

      return c.body(null, StatusCodes.NO_CONTENT);
    },
  );
}
