import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import {
  createReview,
  MovieNotFoundError,
  PlatformNotFoundError,
  ReviewAlreadyExistsError,
} from "../../services/create-review/create-review.service";

const bodySchema = z
  .object({
    rating: z.coerce
      .number()
      .min(0.1)
      .max(5.0)
      .transform((v) => v.toFixed(1))
      .optional(),
    comment: z.string().min(1).optional(),
    watchHistory: z
      .object({
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
      })
      .optional()
      .nullable(),
  })
  .refine((data) => data.rating !== undefined || data.comment !== undefined, {
    message: "rating or comment is required",
    path: ["rating"],
  });

const paramSchema = z.object({
  id: z.coerce.number().int().positive(),
});

export default function (app: AppType) {
  app.post(
    "/movies/:id/reviews",

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
      const { id: movieId } = c.req.valid("param");
      const body = c.req.valid("json");

      const result = await createReview(logger, {
        userId: operator.id,
        movieId,
        rating: body.rating ?? null,
        comment: body.comment ?? null,
        watchHistory: body.watchHistory
          ? {
              platformId: body.watchHistory.platformId,
              watchedDate: body.watchHistory.watchedDate ?? null,
            }
          : null,
      });

      if (result.isErr()) {
        const error = result.error;
        if (error instanceof MovieNotFoundError) {
          return c.json({ error: "movie not found" }, StatusCodes.NOT_FOUND);
        }
        if (error instanceof ReviewAlreadyExistsError) {
          return c.json(
            { error: "review already exists" },
            StatusCodes.CONFLICT,
          );
        }
        if (error instanceof PlatformNotFoundError) {
          return c.json({ error: "platform not found" }, StatusCodes.NOT_FOUND);
        }
      }

      return c.body(null, StatusCodes.CREATED);
    },
  );
}
