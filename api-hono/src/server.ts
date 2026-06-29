import { app } from "./core/app";
import { environment } from "./core/environment";

import "@/app/auth";
import "@/app/health";
import "@/app/master";
import "@/app/movie";
import "@/app/review";
import "@/app/search";
import "@/app/trending";
import "@/app/user";
import "@/app/watchlist";

export default {
  port: environment.SERVER.PORT,
  fetch: app.fetch,
};
