import { AppType } from "@/core/app";
import { requireAuthMiddleware } from "@/core/middlewares/auth.middleware";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import z from "zod";
import {
  ReviewNotFoundError,
  updateReview,
} from "../../services/update-review/update-review.service";

const paramSchema = z.object({
  id: z.coerce.number().int().positive(),
});

const bodySchema = z
  .object({
    rating: z.coerce
      .number()
      .min(0.1)
      .max(5.0)
      .transform((v) => v.toFixed(1))
      .optional(),
    comment: z.string().min(1).optional(),
  })
  .refine((data) => data.rating !== undefined || data.comment !== undefined, {
    message: "rating or comment is required",
    path: ["rating"],
  });

export default function (app: AppType) {
  app.put(
    "/reviews/:id",

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

      const result = await updateReview(
        logger,
        operator.id,
        reviewId,
        body.rating ?? null,
        body.comment ?? null,
      );

      if (result.isErr()) {
        if (result.error instanceof ReviewNotFoundError) {
          return c.json({ error: "review not found" }, StatusCodes.NOT_FOUND);
        }
      }

      return c.body(null, StatusCodes.NO_CONTENT);
    },
  );
}
