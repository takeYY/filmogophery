import { AppType } from "@/core/app";
import { MovieIsNotFound } from "@/core/errors";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import { match, P } from "ts-pattern";
import z from "zod";
import {
  getMovieById,
  MovieDetail,
} from "../../services/movies/movies.service";

const pathSchema = z.object({
  id: z.coerce.number().int().min(1),
});

export default function (app: AppType) {
  app.get(
    `/movies/:id`,

    validator("param", (value, c) => {
      const parsed = pathSchema.safeParse(value);
      if (!parsed.success) {
        return c.json({ error: parsed.error }, 400);
      }
      return parsed.data;
    }),

    requireAuthMiddleware,

    async (c) => {
      const logger = c.get("logger");
      const operator = c.get("operator");
      const path = c.req.valid("param");

      if (!operator) {
        return c.json({ error: "Unauthorized" }, StatusCodes.UNAUTHORIZED);
      }

      const result = await getMovieById(logger, operator, path.id);
      if (result.isOk()) {
        return c.json(result.value satisfies MovieDetail, StatusCodes.OK);
      }

      return match(result.error)
        .with(P.instanceOf(MovieIsNotFound), () =>
          c.json(
            { message: "movie not found", errors: { id: ["218"] } },
            StatusCodes.NOT_FOUND,
          ),
        )
        .exhaustive();
    },
  );
}
