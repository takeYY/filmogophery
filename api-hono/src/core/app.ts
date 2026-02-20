import { Hono } from "hono";

export const app = new Hono().basePath("/v1");

export type AppType = typeof app;
