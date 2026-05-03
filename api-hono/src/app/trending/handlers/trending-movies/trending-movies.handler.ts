import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { TrendingMovie } from "@/core/types/movie";
import { StatusCodes } from "http-status-codes";
import { getTrendingMovies } from "../../services/trending-movies/trending-movies.service";

export default function (app: AppType) {
  app.get(
    "/trending/movies",

    requireAuthMiddleware,

    async (c) => {
      const logger = c.get("logger");
      const operator = c.get("operator")!;

      const result = await getTrendingMovies(logger, operator.id);
      if (result.isErr()) {
        return c.json(
          { error: "failed to fetch trending movies" },
          StatusCodes.INTERNAL_SERVER_ERROR,
        );
      }

      return c.json(result.value satisfies TrendingMovie[], StatusCodes.OK);
    },
  );
}
