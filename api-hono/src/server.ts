import { app } from "./core/app";

import "@/app/health";
import "@/app/movie";
import "@/app/search";

export default {
  port: process.env.SERVER_PORT,
  fetch: app.fetch,
};
