import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { StatusCodes } from "http-status-codes";

export default function (app: AppType) {
  app.get(
    "/users/me",

    requireAuthMiddleware,

    (c) => {
      const operator = c.get("operator")!;

      return c.json(
        { id: operator.id, username: operator.username },
        StatusCodes.OK,
      );
    },
  );
}
