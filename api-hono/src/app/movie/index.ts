import { app } from "@/core/app";

(await import("./handlers/movies/movies.handler")).default(app);

(await import("./handlers/movies-detail/movies-detail.handler")).default(app);

(
  await import("./handlers/movie-watch-history/movie-watch-history.handler")
).default(app);
