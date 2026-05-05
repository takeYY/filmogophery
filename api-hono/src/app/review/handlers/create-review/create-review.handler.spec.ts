import { Variables } from "@/core/app";
import * as authMiddleware from "@/core/middlewares/auth.middleware";
import { afterEach, beforeEach, describe, expect, spyOn, test } from "bun:test";
import { Hono } from "hono";
import { StatusCodes } from "http-status-codes";
import { err, ok } from "neverthrow";
import pino from "pino";
import * as createReviewService from "../../services/create-review/create-review.service";
import createReviewHandler from "./create-review.handler";

const testLogger = pino({ level: "silent" });
const mockUser = { id: 1, name: "Test User", email: "test@example.com" };

describe("create-review.handler", () => {
  beforeEach(() => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockImplementation(
      async (c, next) => {
        c.set("operator", mockUser as any);
        await next();
      },
    );
  });

  afterEach(() => {
    spyOn(createReviewService, "createReview").mockRestore();
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();
  });

  const makeApp = () => {
    const testApp = new Hono<{ Variables: Variables }>().basePath("/v1");
    testApp.use(async (c, next) => {
      c.set("logger", testLogger);
      await next();
    });
    createReviewHandler(testApp);
    return testApp;
  };

  test("should create review with rating only", async () => {
    spyOn(createReviewService, "createReview").mockResolvedValue(ok(undefined));

    const res = await makeApp().request("/v1/movies/1/reviews", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ rating: 4.5 }),
    });

    expect(res.status).toBe(StatusCodes.CREATED);
  });

  test("should create review with comment only", async () => {
    spyOn(createReviewService, "createReview").mockResolvedValue(ok(undefined));

    const res = await makeApp().request("/v1/movies/1/reviews", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ comment: "面白かった" }),
    });

    expect(res.status).toBe(StatusCodes.CREATED);
  });

  test("should create review with watch history", async () => {
    spyOn(createReviewService, "createReview").mockResolvedValue(ok(undefined));

    const res = await makeApp().request("/v1/movies/1/reviews", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        rating: 3.0,
        comment: "良かった",
        watchHistory: { platformId: 1, watchedDate: "2025-01-01" },
      }),
    });

    expect(res.status).toBe(StatusCodes.CREATED);
  });

  test("should return 400 when neither rating nor comment is provided", async () => {
    const res = await makeApp().request("/v1/movies/1/reviews", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({}),
    });

    expect(res.status).toBe(StatusCodes.BAD_REQUEST);
  });

  test("should return 404 when movie not found", async () => {
    spyOn(createReviewService, "createReview").mockResolvedValue(
      err(new createReviewService.MovieNotFoundError()),
    );

    const res = await makeApp().request("/v1/movies/999/reviews", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ rating: 3.0 }),
    });

    expect(res.status).toBe(StatusCodes.NOT_FOUND);
  });

  test("should return 409 when review already exists", async () => {
    spyOn(createReviewService, "createReview").mockResolvedValue(
      err(new createReviewService.ReviewAlreadyExistsError()),
    );

    const res = await makeApp().request("/v1/movies/1/reviews", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ rating: 3.0 }),
    });

    expect(res.status).toBe(StatusCodes.CONFLICT);
  });

  test("should return 401 when not authenticated", async () => {
    spyOn(authMiddleware, "requireAuthMiddleware").mockRestore();

    const res = await makeApp().request("/v1/movies/1/reviews", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ rating: 3.0 }),
    });

    expect(res.status).toBe(StatusCodes.UNAUTHORIZED);
  });
});
