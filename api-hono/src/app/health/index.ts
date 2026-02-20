import { app } from "@/core/app";

(await import("./handlers/check-health.handler")).default(app);
