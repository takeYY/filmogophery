import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { Movie } from "@/core/types/movie";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import { getMovies } from "../../services/movies/movies.service";

const querySchema = z.object({
  genre: z.string().optional(),
  limit: z.coerce.number().int().positive().default(12),
  offset: z.coerce.number().int().min(0).default(0),
});

export default function (app: AppType) {
  app.get(
    `/movies`,
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

      const result = await getMovies(
        logger,
        operator.id,
        query.genre,
        query.limit,
        query.offset,
      );
      return c.json(result.value satisfies Movie[], StatusCodes.OK);
    },
  );
}
