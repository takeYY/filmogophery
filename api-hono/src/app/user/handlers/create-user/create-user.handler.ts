import { AppType } from "@/core/app";
import { validator } from "hono/validator";
import { StatusCodes } from "http-status-codes";
import { match, P } from "ts-pattern";
import z from "zod";
import {
  createUser,
  UserAlreadyExistsError,
} from "../../services/create-user/create-user.service";

const passwordSchema = z
  .string()
  .min(8)
  .refine((v) => /[A-Z]/.test(v) && /[a-z]/.test(v) && /[0-9]/.test(v), {
    message: "password must contain uppercase, lowercase, and digit",
  });

const bodySchema = z.object({
  username: z.string().min(1),
  email: z.email(),
  password: passwordSchema,
});

export default function (app: AppType) {
  app.post(
    `/users`,
    validator("json", (value, c) => {
      const parsed = bodySchema.safeParse(value);
      if (!parsed.success) {
        return c.json({ error: parsed.error }, StatusCodes.BAD_REQUEST);
      }
      return parsed.data;
    }),

    async (c) => {
      const logger = c.get("logger");
      const body = c.req.valid("json");

      const result = await createUser(
        logger,
        body.username,
        body.email,
        body.password,
      );

      if (result.isOk()) {
        return c.json(result.value, StatusCodes.CREATED);
      }

      return match(result.error)
        .with(P.instanceOf(UserAlreadyExistsError), () =>
          c.json(
            {
              message: "conflict",
              errors: { username: ["username or email is already taken"] },
            },
            StatusCodes.CONFLICT,
          ),
        )
        .otherwise(() =>
          c.json(
            { message: "system error" },
            StatusCodes.INTERNAL_SERVER_ERROR,
          ),
        );
    },
  );
}
