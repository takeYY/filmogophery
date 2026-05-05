import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import pino from "pino";
import * as updateReviewService from "../../services/update-review/update-review.service";
import updateReviewHandler from "./update-review.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

describe("update-review.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(updateReviewService, "updateReview").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    updateReviewHandler(testApp);
    return testApp;
  };

  const putReview = (
    app: ReturnType<typeof makeApp>,
    id: number,
    body: unknown,
  ) =>
    app.request(`/v1/reviews/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

  test("should update review with rating only", async () => {
    spyOn(updateReviewService, "updateReview").mockResolvedValue(ok(undefined));

    const res = await putReview(makeApp(), 1, { rating: 4.5 });

    expect(res.status).toBe(StatusCodes.NO_CONTENT);
  });

  test("should update review with comment only", async () => {
    spyOn(updateReviewService, "updateReview").mockResolvedValue(ok(undefined));

    const res = await putReview(makeApp(), 1, { comment: "面白かった" });

    expect(res.status).toBe(StatusCodes.NO_CONTENT);
  });

  test("should update review with both rating and comment", async () => {
    spyOn(updateReviewService, "updateReview").mockResolvedValue(ok(undefined));

    const res = await putReview(makeApp(), 1, {
      rating: 3.0,
      comment: "良かった",
    });

    expect(res.status).toBe(StatusCodes.NO_CONTENT);
  });

  test("should return 400 when neither rating nor comment is provided", async () => {
    const res = await putReview(makeApp(), 1, {});

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 400 when rating is out of range", async () => {
    const res = await putReview(makeApp(), 1, { rating: 5.1 });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 404 when review not found", async () => {
    spyOn(updateReviewService, "updateReview").mockResolvedValue(
      err(new updateReviewService.ReviewNotFoundError()),
    );

    const res = await putReview(makeApp(), 999, { rating: 3.0 });

    expect(res.status).toBe(StatusCodes.NOT_FOUND);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await putReview(makeApp(), 1, { rating: 3.0 });

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
