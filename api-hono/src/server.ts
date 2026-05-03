import { app } from "./core/app";

import "@/app/auth";
import "@/app/health";
import "@/app/movie";
import "@/app/search";
import "@/app/trending";
import "@/app/user";

export default {
  port: process.env.SERVER_PORT,
  fetch: app.fetch,
};
