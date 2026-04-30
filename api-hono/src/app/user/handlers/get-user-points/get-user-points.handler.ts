import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import { getUserPoints } from "../../services/get-user-points/get-user-points.service";

const querySchema = z.object({
  limit: z.number().int().min(1).max(50).default(20),
  offset: z.number().int().min(0).default(0),
});

export default function (app: AppType) {
  app.get(
    `/users/me/points`,
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
      const operator = c.get("operator");
      const query = c.req.valid("query");

      if (!operator) {
        return c.json({ error: "Unauthorized" }, StatusCodes.UNAUTHORIZED);
      }

      const result = await getUserPoints(
        logger,
        operator.id,
        query.limit,
        query.offset,
      );

      return c.json(result.value, StatusCodes.OK);
    },
  );
}
