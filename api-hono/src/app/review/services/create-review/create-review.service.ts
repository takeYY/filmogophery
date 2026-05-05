import { err, ok, Result } from "neverthrow";
import { Logger } from "pino";
import {
  createReviewWithWatchHistory,
  fetchMovieById,
  fetchPlatformById,
  fetchReviewByMovieId,
} from "../../repositories/review.repository";

export class MovieNotFoundError extends Error {}
export class PlatformNotFoundError extends Error {}
export class ReviewAlreadyExistsError extends Error {}

export type WatchHistoryInput = {
  platformId: number;
  watchedDate: string | null;
};

export type CreateReviewInput = {
  userId: number;
  movieId: number;
  rating: string | null;
  comment: string | null;
  watchHistory: WatchHistoryInput | null;
};

export async function createReview(
  logger: Logger,
  input: CreateReviewInput,
): Promise<
  Result<
    void,
    MovieNotFoundError | PlatformNotFoundError | ReviewAlreadyExistsError
  >
> {
  logger.info(
    { movieId: input.movieId, userId: input.userId },
    "createReview called",
  );

  // 映画の存在確認
  const movie = await fetchMovieById(input.movieId);
  if (!movie) {
    logger.info({ movieId: input.movieId }, "movie not found");
    return err(new MovieNotFoundError());
  }

  // レビュー重複確認
  const existing = await fetchReviewByMovieId(input.userId, input.movieId);
  if (existing) {
    logger.info(
      { movieId: input.movieId, userId: input.userId },
      "review already exists",
    );
    return err(new ReviewAlreadyExistsError());
  }

  // プラットフォームの存在確認（視聴履歴あり時のみ）
  if (input.watchHistory !== null) {
    const platform = await fetchPlatformById(input.watchHistory.platformId);
    if (!platform) {
      logger.info(
        { platformId: input.watchHistory.platformId },
        "platform not found",
      );
      return err(new PlatformNotFoundError());
    }
  }

  // トランザクションでレビュー・視聴履歴・ポイントを一括登録
  await createReviewWithWatchHistory(
    {
      userId: input.userId,
      movieId: input.movieId,
      rating: input.rating,
      comment: input.comment,
    },
    input.watchHistory
      ? {
          userId: input.userId,
          movieId: input.movieId,
          platformId: input.watchHistory.platformId,
          watchedDate: input.watchHistory.watchedDate,
        }
      : null,
    movie.runtimeMinutes,
  );

  logger.info(
    { movieId: input.movieId, userId: input.userId },
    "createReview completed",
  );
  return ok(undefined);
}
