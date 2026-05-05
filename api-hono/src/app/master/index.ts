import { app } from "@/core/app";

(await import("./handlers/get-genres/get-genres.handler")).default(app);
