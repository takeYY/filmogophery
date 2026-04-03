import { dbConnections } from "@/core/db";
import { reviews } from "@/core/drizzle/schema";
import { and, eq } from "drizzle-orm";
import { MySql2Database } from "drizzle-orm/mysql2";

export async function fetchReviewByMovieId(
  userId: number,
  movieId: number,
  db: MySql2Database = dbConnections.readonly,
) {
  return (
    await db
      .select({
        id: reviews.id,
        createdAt: reviews.createdAt,
        updatedAt: reviews.updatedAt,
        rating: reviews.rating,
        comment: reviews.comment,
      })
      .from(reviews)
      .where(and(eq(reviews.userId, userId), eq(reviews.movieId, movieId)))
      .limit(1)
  ).at(0);
}
