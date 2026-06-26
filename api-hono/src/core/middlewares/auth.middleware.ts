import { eq } from "drizzle-orm";
import { createMiddleware } from "hono/factory";
import { verify } from "hono/jwt";
import { StatusCodes } from "http-status-codes";
import { Variables } from "../app";
import { dbConnections } from "../db";
import { users } from "../drizzle/schema";
import { environment } from "../environment";

type JWTPayload = {
  user_id: number;
  exp?: number;
};

export const requireAuthMiddleware = createMiddleware<{
  Variables: Variables;
}>(async (c, next) => {
  const logger = c.get("logger");

  // Authorizationヘッダーからトークンを取得
  const authHeader = c.req.header("Authorization");
  if (!authHeader) {
    return c.json({ error: "Unauthorized" }, StatusCodes.UNAUTHORIZED);
  }

  // Bearer プレフィックスを削除
  if (!authHeader.startsWith("Bearer ")) {
    return c.json({ error: "Unauthorized" }, StatusCodes.UNAUTHORIZED);
  }
  const tokenString = authHeader.slice(7);

  // JWT トークンを検証
  let payload: JWTPayload;
  try {
    payload = (await verify(
      tokenString,
      environment.TOKEN.JWT_SECRET,
      "HS256",
    )) as JWTPayload;
  } catch (err) {
    logger.error({ err }, "invalid token");
    return c.json({ error: "Unauthorized" }, StatusCodes.UNAUTHORIZED);
  }

  // ユーザーをDBから取得
  const [user] = await dbConnections.readonly
    .select()
    .from(users)
    .where(eq(users.id, payload.user_id))
    .limit(1);

  if (!user) {
    logger.error({ userId: payload.user_id }, "user is not found");
    return c.json({ error: "Unauthorized" }, StatusCodes.UNAUTHORIZED);
  }

  c.set("operator", user);
  await next();
});
