import { app } from "@/core/app";

(await import("./handlers/create-review/create-review.handler")).default(app);
(await import("./handlers/update-review/update-review.handler")).default(app);
(
  await import("./handlers/create-watch-history/create-watch-history.handler")
).default(app);
