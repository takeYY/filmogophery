import { app } from "@/core/app";

(await import("./handlers/get-watchlist/get-watchlist.handler")).default(app);
(await import("./handlers/add-watchlist/add-watchlist.handler")).default(app);
(await import("./handlers/delete-watchlist/delete-watchlist.handler")).default(
  app,
);
