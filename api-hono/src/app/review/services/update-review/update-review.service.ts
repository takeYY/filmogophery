import { err, ok, Result } from "neverthrow";
import { Logger } from "pino";
import {
  fetchReviewById,
  updateReview as updateReviewInDB,
} from "../../repositories/review.repository";

export class ReviewNotFoundError extends Error {}

export async function updateReview(
  logger: Logger,
  userId: number,
  reviewId: number,
  rating: string | null,
  comment: string | null,
): Promise<Result<void, ReviewNotFoundError>> {
  logger.info({ userId, reviewId }, "updateReview called");

  // レビューの存在確認（ユーザー所有チェック込み）
  const review = await fetchReviewById(userId, reviewId);
  if (!review) {
    logger.info({ reviewId }, "review not found");
    return err(new ReviewNotFoundError());
  }

  await updateReviewInDB(reviewId, rating, comment);

  logger.info({ reviewId }, "updateReview completed");
  return ok(undefined);
}
