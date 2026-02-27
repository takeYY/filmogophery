import { app } from "./core/app";

import "@/app/health";
import "@/app/movie";

export default {
  port: process.env.SERVER_PORT,
  fetch: app.fetch,
};
