import { app } from "./core/app";

import "@/app/auth";
import "@/app/health";
import "@/app/movie";
import "@/app/review";
import "@/app/search";
import "@/app/trending";
import "@/app/user";
import "@/app/watchlist";

export default {
  port: process.env.SERVER_PORT,
  fetch: app.fetch,
};
