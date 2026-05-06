import { defaultRuntimeMinutes } from "@/core/definition";
import { err, ok, Result } from "neverthrow";
import { Logger } from "pino";
import {
  createWatchHistoryWithPoints,
  fetchMovieById,
  fetchPlatformById,
  fetchReviewById,
} from "../../repositories/review.repository";

export class ReviewNotFoundError extends Error {}
export class PlatformNotFoundError extends Error {}

export type CreateWatchHistoryInput = {
  userId: number;
  reviewId: number;
  platformId: number;
  watchedDate: string | null;
};

export async function createWatchHistory(
  logger: Logger,
  input: CreateWatchHistoryInput,
): Promise<Result<void, ReviewNotFoundError | PlatformNotFoundError>> {
  logger.info(
    { reviewId: input.reviewId, userId: input.userId },
    "createWatchHistory called",
  );

  // レビューの存在確認（ユーザー所有チェック込み）
  const review = await fetchReviewById(input.userId, input.reviewId);
  if (!review) {
    logger.info({ reviewId: input.reviewId }, "review not found");
    return err(new ReviewNotFoundError());
  }

  // プラットフォームの存在確認
  const platform = await fetchPlatformById(input.platformId);
  if (!platform) {
    logger.info({ platformId: input.platformId }, "platform not found");
    return err(new PlatformNotFoundError());
  }

  // 映画情報を取得（ポイント計算に上映時間が必要）
  const movie = await fetchMovieById(review.movieId);
  const runtimeMinutes = movie?.runtimeMinutes ?? defaultRuntimeMinutes;

  await createWatchHistoryWithPoints(
    {
      userId: input.userId,
      movieId: review.movieId,
      platformId: input.platformId,
      watchedDate: input.watchedDate,
    },
    runtimeMinutes,
  );

  logger.info({ reviewId: input.reviewId }, "createWatchHistory completed");
  return ok(undefined);
}
