import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import {
  addWatchlist,
  MovieNotFoundError,
} from "../../services/add-watchlist/add-watchlist.service";

const bodySchema = z.object({
  movieId: z.number().int().positive(),
  priority: z.number().int().min(1).max(5).default(1),
});

export default function (app: AppType) {
  app.post(
    "/watchlist",

    validator("json", (value, c) => {
      const parsed = bodySchema.safeParse(value);
      if (!parsed.success) {
        return c.json({ error: parsed.error }, StatusCodes.BAD_REQUEST);
      }
      return parsed.data;
    }),

    requireAuthMiddleware,

    async (c) => {
      const logger = c.get("logger");
      const operator = c.get("operator")!;
      const body = c.req.valid("json");

      const result = await addWatchlist(
        logger,
        operator.id,
        body.movieId,
        body.priority,
      );

      if (result.isErr()) {
        if (result.error instanceof MovieNotFoundError) {
          return c.json({ error: "movie not found" }, StatusCodes.NOT_FOUND);
        }
      }

      return c.body(null, StatusCodes.CREATED);
    },
  );
}
