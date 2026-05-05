import { app } from "@/core/app";

(await import("./handlers/create-review/create-review.handler")).default(app);
