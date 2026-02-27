import { app } from "@/core/app";

(await import("./handlers/movies/movies.handler")).default(app);
