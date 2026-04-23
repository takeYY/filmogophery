import { app } from "./core/app";

import "@/app/health";
import "@/app/movie";
import "@/app/search";
import "@/app/user";

export default {
  port: process.env.SERVER_PORT,
  fetch: app.fetch,
};
