import { dbConnections } from "@/core/db";
import { refreshTokens } from "@/core/drizzle/schema";
import { environment } from "@/core/environment";
import crypto from "crypto";
import { sign } from "hono/jwt";
import { err, ok } from "neverthrow";
import { Logger } from "pino";
import { insertUser as _insertUser } from "../../repositories/users/users.repository";

const EXPIRES_IN = 3600; // 1時間

export class UserAlreadyExistsError extends Error {}

export type Token = {
  accessToken: string;
  refreshToken: string;
  tokenType: string;
  expiresIn: number;
  expiresAt: string;
};

type Deps = {
  insertUser?: typeof _insertUser;
};

export async function createUser(
  logger: Logger,
  username: string,
  email: string,
  password: string,
  { insertUser = _insertUser }: Deps = {},
) {
  logger.info({ username, email }, "createUser called");

  const now = new Date();
  const expiresAt = new Date(now.getTime() + EXPIRES_IN * 1000);

  // パスワードをハッシュ化
  const passwordHash = await Bun.password.hash(password, {
    algorithm: "bcrypt",
  });

  const db = dbConnections.default;

  try {
    const userId = await insertUser({
      username,
      email,
      passwordHash,
      lastLoginAt: now.toISOString().slice(0, 19).replace("T", " "),
    });

    // アクセストークンを生成
    const accessToken = await sign(
      { user_id: userId, exp: Math.floor(expiresAt.getTime() / 1000) },
      environment.TOKEN.JWT_SECRET!,
      "HS256",
    );

    // リフレッシュトークンを生成してDBに保存
    const refreshToken = crypto.randomBytes(32).toString("hex");
    const tokenHash = crypto
      .createHash("sha256")
      .update(refreshToken)
      .digest("hex");
    const refreshExpiresAt = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000);

    await db.insert(refreshTokens).values({
      userId,
      tokenHash,
      expiresAt: refreshExpiresAt.toISOString().slice(0, 19).replace("T", " "),
    });

    logger.info({ userId }, "successfully created user");

    return ok({
      accessToken,
      refreshToken,
      tokenType: "Bearer",
      expiresIn: EXPIRES_IN,
      expiresAt: expiresAt.toISOString(),
    } satisfies Token);
  } catch (e: unknown) {
    const msg = e instanceof Error ? e.message : String(e);
    if (msg.includes("Duplicate entry") || msg.includes("duplicated key")) {
      return err(
        new UserAlreadyExistsError("username or email is already taken"),
      );
    }
    logger.error({ err: e }, "failed to create user");
    return err(e as Error);
  }
}
