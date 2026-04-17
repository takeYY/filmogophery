import { app } from "@/core/app";

(await import("./handlers/search-movies/search-movies.handler")).default(app);
