import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { StatusCodes } from "http-status-codes";
import { revokeActiveTokensByUserId } from "../../repositories/users/users.repository";

export default function (app: AppType) {
  app.post(
    "/auth/logout",

    requireAuthMiddleware,

    async (c) => {
      const logger = c.get("logger");
      const operator = c.get("operator")!;

      const now = new Date().toISOString().replace("T", " ").replace("Z", "");
      await revokeActiveTokensByUserId(operator.id, now);

      logger.info({ userId: operator.id }, "logout succeeded");
      return c.body(null, StatusCodes.NO_CONTENT);
    },
  );
}
