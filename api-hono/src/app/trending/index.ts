import { app } from "@/core/app";

(await import("./handlers/trending-movies/trending-movies.handler")).default(
  app,
);
