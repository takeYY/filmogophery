import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { Movie } from "@/core/types/movie";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import { match, P } from "ts-pattern";
import z from "zod";
import { searchMovies } from "../../services/search-movies/search-movies.service";

const querySchema = z.object({
  title: z.string(),
  limit: z.number().int().min(1).max(12).default(12),
  offset: z.number().int().min(0).default(0),
});

export default function (app: AppType) {
  app.get(
    `/search/movies`,
    validator("query", (value, c) => {
      const parsed = querySchema.safeParse(value);
      if (!parsed.success) {
        return c.json({ error: parsed.error }, 400);
      }
      return parsed.data;
    }),

    requireAuthMiddleware,

    async (c) => {
      const logger = c.get("logger");
      const operator = c.get("operator");
      const query = c.req.valid("query");

      if (!operator) {
        return c.json({ error: "Unauthorized" }, StatusCodes.UNAUTHORIZED);
      }

      const result = await searchMovies(
        logger,
        query.title,
        query.limit,
        query.offset,
      );
      if (result.isOk()) {
        return c.json(result.value satisfies Movie[], StatusCodes.OK);
      }

      return match(result.error)
        .with(P.instanceOf(Error), () =>
          c.json(
            { message: "system error", errors: { title: query.title } },
            StatusCodes.INTERNAL_SERVER_ERROR,
          ),
        )
        .exhaustive();
    },
  );
}
