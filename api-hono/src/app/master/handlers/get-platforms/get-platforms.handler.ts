import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { StatusCodes } from "http-status-codes";
import { fetchAllPlatforms } from "../../repositories/platforms.repository";

export default function (app: AppType) {
  app.get(
    "/platforms",

    requireAuthMiddleware,

    async (c) => {
      const platforms = await fetchAllPlatforms();
      return c.json(platforms, StatusCodes.OK);
    },
  );
}
