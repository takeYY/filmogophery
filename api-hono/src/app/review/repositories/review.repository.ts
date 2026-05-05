import { dbConnections } from "@/core/db";
import {
  movies,
  platforms,
  pointHistory,
  reviews,
  userPoints,
  watchHistory,
} from "@/core/drizzle/schema";
import { and, eq, sql } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

/**
 * 映画IDに一致する映画を取得する
 */
export async function fetchMovieById(
  movieId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const [movie] = await db
    .select({ id: movies.id, runtimeMinutes: movies.runtimeMinutes })
    .from(movies)
    .where(eq(movies.id, movieId))
    .limit(1);
  return movie ?? null;
}

/**
 * プラットフォームIDに一致するプラットフォームを取得する
 */
export async function fetchPlatformById(
  platformId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const [platform] = await db
    .select({ id: platforms.id })
    .from(platforms)
    .where(eq(platforms.id, platformId))
    .limit(1);
  return platform ?? null;
}

/**
 * ユーザーの映画レビューが既に存在するか確認する
 */
export async function fetchReviewByMovieId(
  userId: number,
  movieId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const [review] = await db
    .select({ id: reviews.id })
    .from(reviews)
    .where(and(eq(reviews.userId, userId), eq(reviews.movieId, movieId)))
    .limit(1);
  return review ?? null;
}

/**
 * レビューIDとユーザーIDに一致するレビューを取得する
 */
export async function fetchReviewById(
  userId: number,
  reviewId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  const [review] = await db
    .select({
      id: reviews.id,
      userId: reviews.userId,
      movieId: reviews.movieId,
    })
    .from(reviews)
    .where(and(eq(reviews.id, reviewId), eq(reviews.userId, userId)))
    .limit(1);
  return review ?? null;
}

/**
 * レビューを更新する
 */
export async function updateReview(
  reviewId: number,
  rating: string | null,
  comment: string | null,
  db: MySql2Database = dbConnections.default,
) {
  await db
    .update(reviews)
    .set({
      rating: rating ?? undefined,
      comment: comment ?? undefined,
    })
    .where(eq(reviews.id, reviewId));
}

export type CreateReviewInput = {
  userId: number;
  movieId: number;
  rating: string | null;
  comment: string | null;
};

export type CreateWatchHistoryInput = {
  userId: number;
  movieId: number;
  platformId: number;
  watchedDate: string | null;
};

/**
 * レビュー・視聴履歴・ポイントをトランザクションで一括登録する
 */
export async function createReviewWithWatchHistory(
  reviewInput: CreateReviewInput,
  watchHistoryInput: CreateWatchHistoryInput | null,
  runtimeMinutes: number,
  db: MySql2Database = dbConnections.default,
) {
  return db.transaction(async (tx) => {
    // レビューを登録
    const [reviewResult] = await tx.insert(reviews).values({
      userId: reviewInput.userId,
      movieId: reviewInput.movieId,
      rating: reviewInput.rating ?? undefined,
      comment: reviewInput.comment ?? undefined,
    });
    const reviewId = reviewResult.insertId;

    // レビューポイントを付与（20pt固定）
    const REVIEW_POINTS = 20;
    await grantPoints(
      tx,
      reviewInput.userId,
      REVIEW_POINTS,
      "review",
      reviewId,
    );

    // 視聴履歴を登録（入力がある場合のみ）
    if (watchHistoryInput !== null) {
      const [whResult] = await tx.insert(watchHistory).values({
        userId: watchHistoryInput.userId,
        movieId: watchHistoryInput.movieId,
        platformId: watchHistoryInput.platformId,
        watchedDate: watchHistoryInput.watchedDate ?? undefined,
      });
      const watchHistoryId = whResult.insertId;

      // 視聴履歴ポイントを付与（上映時間による段階的ポイント）
      const watchPoints = calcWatchPoints(runtimeMinutes);
      await grantPoints(
        tx,
        watchHistoryInput.userId,
        watchPoints,
        "watch_history",
        watchHistoryId,
      );
    }
  });
}

/**
 * 視聴履歴・ポイントをトランザクションで登録する
 */
export async function createWatchHistoryWithPoints(
  input: CreateWatchHistoryInput,
  runtimeMinutes: number,
  db: MySql2Database = dbConnections.default,
) {
  return db.transaction(async (tx) => {
    const [whResult] = await tx.insert(watchHistory).values({
      userId: input.userId,
      movieId: input.movieId,
      platformId: input.platformId,
      watchedDate: input.watchedDate ?? undefined,
    });
    const watchHistoryId = whResult.insertId;

    const watchPoints = calcWatchPoints(runtimeMinutes);
    await grantPoints(
      tx,
      input.userId,
      watchPoints,
      "watch_history",
      watchHistoryId,
    );
  });
}

/**
 * ポイントを付与してポイント履歴を記録する
 */
async function grantPoints(
  tx: MySql2Database,
  userId: number,
  points: number,
  action: string,
  referenceId: number,
) {
  // user_points を upsert（なければ作成、あれば加算）
  await tx
    .insert(userPoints)
    .values({ userId, totalPoints: points, level: 1 })
    .onDuplicateKeyUpdate({
      set: { totalPoints: sql`${userPoints.totalPoints} + ${points}` },
    });

  // ポイント履歴を記録
  await tx.insert(pointHistory).values({ userId, points, action, referenceId });
}

/**
 * 上映時間からポイントを計算する（Echoと同じロジック）
 */
function calcWatchPoints(runtimeMinutes: number): number {
  if (runtimeMinutes <= 90) return 10;
  if (runtimeMinutes <= 150) return 15;
  return 20;
}
