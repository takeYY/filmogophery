import { app } from "@/core/app";

(await import("./handlers/create-user/create-user.handler")).default(app);
(await import("./handlers/get-current-user/get-current-user.handler")).default(
  app,
);
(await import("./handlers/get-user-points/get-user-points.handler")).default(
  app,
);
(
  await import("./handlers/get-watch-history/get-watch-history.handler")
).default(app);
