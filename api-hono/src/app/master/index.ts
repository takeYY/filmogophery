import { app } from "@/core/app";

(await import("./handlers/get-genres/get-genres.handler")).default(app);
(await import("./handlers/get-platforms/get-platforms.handler")).default(app);
