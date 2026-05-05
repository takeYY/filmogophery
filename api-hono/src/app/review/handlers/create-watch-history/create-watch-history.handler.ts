import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import {
  createWatchHistory,
  PlatformNotFoundError,
  ReviewNotFoundError,
} from "../../services/create-watch-history/create-watch-history.service";

const paramSchema = z.object({
  id: z.coerce.number().int().positive(),
});

const bodySchema = z.object({
  platformId: z.number().int().positive(),
  watchedDate: z
    .string()
    .regex(/^\d{4}-\d{2}-\d{2}$/, "watchedDate must be YYYY-MM-DD")
    .refine(
      (d) => new Date(d) <= new Date(),
      "watchedDate cannot be in the future",
    )
    .optional()
    .nullable(),
});

export default function (app: AppType) {
  app.post(
    "/reviews/:id/history",

    validator("param", (value, c) => {
      const parsed = paramSchema.safeParse(value);
      if (!parsed.success) {
        return c.json({ error: parsed.error }, StatusCodes.BAD_REQUEST);
      }
      return parsed.data;
    }),

    validator("json", (value, c) => {
      const parsed = bodySchema.safeParse(value);
      if (!parsed.success) {
        return c.json({ error: parsed.error }, StatusCodes.BAD_REQUEST);
      }
      return parsed.data;
    }),

    requireAuthMiddleware,

    async (c) => {
      const logger = c.get("logger");
      const operator = c.get("operator")!;
      const { id: reviewId } = c.req.valid("param");
      const body = c.req.valid("json");

      const result = await createWatchHistory(logger, {
        userId: operator.id,
        reviewId,
        platformId: body.platformId,
        watchedDate: body.watchedDate ?? null,
      });

      if (result.isErr()) {
        if (result.error instanceof ReviewNotFoundError) {
          return c.json({ error: "review not found" }, StatusCodes.NOT_FOUND);
        }
        if (result.error instanceof PlatformNotFoundError) {
          return c.json({ error: "platform not found" }, StatusCodes.NOT_FOUND);
        }
      }

      return c.body(null, StatusCodes.CREATED);
    },
  );
}
