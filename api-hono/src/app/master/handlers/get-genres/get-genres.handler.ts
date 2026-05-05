import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { StatusCodes } from "http-status-codes";
import { fetchAllGenres } from "../../repositories/genres.repository";

export default function (app: AppType) {
  app.get(
    "/genres",

    requireAuthMiddleware,

    async (c) => {
      const genres = await fetchAllGenres();
      return c.json(genres, StatusCodes.OK);
    },
  );
}
