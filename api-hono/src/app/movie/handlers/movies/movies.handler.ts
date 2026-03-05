import { AppType } from "@/core/app";
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

    async (c) => {
      const query = c.req.valid("query");

      const result = await getMovies(query.genre, query.limit, query.offset);
      return c.json(result.value, StatusCodes.OK);
    },
  );
}
