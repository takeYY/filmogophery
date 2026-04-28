import { dbConnections } from "@/core/db";
import { refreshTokens } from "@/core/drizzle/schema";
import { environment } from "@/core/environment";
import crypto from "crypto";
import { sign } from "hono/jwt";
import { err, ok } from "neverthrow";
import { Logger } from "pino";
import {
  findUserByEmail as _findUserByEmail,
  updateLastLoginAt as _updateLastLoginAt,
} from "../../repositories/users/users.repository";

const EXPIRES_IN = 3600; // 1時間

export class InvalidCredentialsError extends Error {}

export type Token = {
  accessToken: string;
  refreshToken: string;
  tokenType: string;
  expiresIn: number;
  expiresAt: string;
};

export async function insertRefreshToken(
  userId: number,
  tokenHash: string,
  expiresAt: string,
) {
  await dbConnections.default.insert(refreshTokens).values({
    userId,
    tokenHash,
    expiresAt,
  });
}

type Deps = {
  findUserByEmail?: typeof _findUserByEmail;
  updateLastLoginAt?: typeof _updateLastLoginAt;
  insertRefreshToken?: typeof insertRefreshToken;
};

export async function login(
  logger: Logger,
  email: string,
  password: string,
  {
    findUserByEmail = _findUserByEmail,
    updateLastLoginAt = _updateLastLoginAt,
    insertRefreshToken: _insertRefreshToken = insertRefreshToken,
  }: Deps = {},
) {
  logger.info({ email }, "login called");

  const user = await findUserByEmail(email);
  if (!user) {
    return err(new InvalidCredentialsError("invalid email or password"));
  }

  const isValid = await Bun.password.verify(password, user.passwordHash);
  if (!isValid) {
    return err(new InvalidCredentialsError("invalid email or password"));
  }

  const now = new Date();
  const expiresAt = new Date(now.getTime() + EXPIRES_IN * 1000);

  try {
    // アクセストークンを生成
    const accessToken = await sign(
      { user_id: user.id, exp: Math.floor(expiresAt.getTime() / 1000) },
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

    await _insertRefreshToken(
      user.id,
      tokenHash,
      refreshExpiresAt.toISOString().slice(0, 19).replace("T", " "),
    );

    await updateLastLoginAt(
      user.id,
      now.toISOString().slice(0, 19).replace("T", " "),
    );

    logger.info({ userId: user.id }, "successfully logged in");

    return ok({
      accessToken,
      refreshToken,
      tokenType: "Bearer",
      expiresIn: EXPIRES_IN,
      expiresAt: expiresAt.toISOString(),
    } satisfies Token);
  } catch (e) {
    logger.error({ err: e }, "failed to login");
    return err(e as Error);
  }
}
