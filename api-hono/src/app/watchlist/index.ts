import { app } from "@/core/app";

(await import("./handlers/get-watchlist/get-watchlist.handler")).default(app);
